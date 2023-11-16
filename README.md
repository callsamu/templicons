# Templicons

This package implements Icon components for the [Templ](https://github.com/tempojs/tempo)
templating language by using the Iconify API, which are pretty similar to the Iconify React,
Vue, and other frameworks components, except that retrieved icons are cached on your Go application, 
which means that only one request per icon is needed and no javascript is required.

## Usage

### Singleton 

Once you import this package, a cache singleton will
be created, making usage pretty straightforward:

```templ
package foobar

import i "github.com/callsamu/templicons"

templ FooComponent() {
    <!-- A Simple Home Icon -->
    @i.Icon("mdi:home", nil)
    <!-- The Same Icon, But With Width and Height Set -->
    @i.Icon("mdi:home", &i.Parameters{}.SetDimensions(20, 20))
    <!-- The Same Icon, But With Color Set to Blue -->
    @i.Icon("mdi:home", &i.Parameters{}.SetColor("#FF0000"))
    <!-- Or Just Ignore Any Setters -->
    @i.Icon("mdi:home", &i.Parameters{Width: 20, Height: 20, Color: "#FF0000"})
}
```

Currently, the "Width", "Height" and "Color" are the only supported parameters. Also,
keep in mind that an invalid icon format will result in a panic.

### Instantiating The Cache

However, if you would like to avoid global state, you can manually create a cache
instance, which supports all of the singleton functions, and pass it to your components.
This approach is a little bit more manual, and requires some configuration.

```templ
package main

import "github.com/callsamu/templicons"

templ FooComponent(c *templicons.Cache) {
    @c.Icon("mdi:home", nil)
}
```

```go
package main

import (
    "os"
    "context"

    "github.com/callsamu/templicons"
)

func main() {
    // The Iconify Instance You Would Like to Use
    instance := "https://api.iconify.design"
    // The Client
    client := templicons.NewIconifyClient()
    // The Fallback Component (we will come back to this later)
    fallback := templicons.DefaultFallback

    // Your cache instance
    cache := templicons.NewCache(api, client, fallback)

    component := FooComponent(cache)
    err := component.Render(context.Background(), os.Stdout)
    if err != nil {
        panic(err)
    }
}
```

### Fallbacks

A major problem with the `Icon` (function/method), is that it is *blocking*:
it will directly attempt to fetch the icon from the API. This may not be 
acceptable, even if all of the results are cached. To solve this, you
can use the `IconWithFallback` (function or method, it depends if you are using
the singleton or not). Instead of waiting for the expensive API call, it
renders a *Fallback Component*, which can be easily customized.

```templ
package foobar

import i "github.com/callsamu/templicons"

templ MyFallback(text string, url string) {
    <span> { text } </span>
}

templ FooComponent() {
    <!-- It requires an additional parameter -->
    <!-- which is a placeholder for using on -->
    <!-- the fallback component -->
    @i.IconWithFallback("mdi:home", "home", nil)
}
```

```go
package main

import (
    "os"
    "log"
    "context"

    "github.com/callsamu/templicons"
)

func main() {
    templicons.SetFallback(MyFallback)

    // Of course you can always use the select statement
    go func() {
        for err := templicons.Errors() {
            log.Error(err)
        }
    }

    component := FooComponent()
    err := component.Render(context.Background(), os.Stdout)
    if err != nil {
        panic(err)
    }
}
```

There are two things to note:

1. Each call of `IconWithFallback` invokes a goroutine which makes
   a API call whose errors are communicated via a channel, which
   be listened.

2. A new fallback component is created and set, overriding the default.
   A fallback component is just any function that takes a text and icon url
   strings as input and returns a templ component.

A nice thing about customizing fallback components is that you can use the url
parameter and a javascript library like HTMX to achieve something like this:

```templ
package foobar

import i "github.com/callsamu/templicons"

templ HTMXFallback(text string, url string) {
	<span 
		hx-get={ url }
		hx-disinherit="*"
		hx-target="this" 
		hx-swap="outerHTML" 
		hx-confirm="unset"
		hx-trigger="load"
		hx-on="htmx:configRequest: event.detail.headers=''">
        { text } 
    </span>
}
```

This fallback component will make a *client side* call to the Iconify API 
and replace itself by the response, thus rendering the icon on the page even 
if it wasn't yet retrieved, while the default fallback would just render
a span containing the text parameters value and left the url untouched.

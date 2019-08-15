## `raffle`

`raffle` is a GopherJS application that generates Twitter-based entries into a simple raffle.

The idea for `raffle` came from [London Gophers](https://gophers.london) meetups. [JetBrains](https://www.jetbrains.com)
kindly support our community by donating a number of 1-year licenses each month. The fairest way to distribute these
amongst meetup attendees is via a simple raffle. Hence `raffle` was born.

During a meetup we announce a "secret" code. A would-be entrant uses this code along with their Twitter handle, and
`raffle` generates a tweet which when posted represents an entry into the raffle.

### Demo

See https://gophers.london/raffle/.

Entering the London Gophers raffle results in a tweet like this:

> Hey @LondonGophers, please enter me into the @jetbrains raffle! bc2e8bba9e5e18c3bc1a7319a0ff04e8358cb7367dd38e029d26f545fdd9fc7a #LondonGophers

The 64-bit hash is the result of hashing the entrant's Twitter handle and the "secret". Entries can then be verified
using the very bare-bones `github.com/LondonGophers/raffle/checker`.

### Using `raffle`

`raffle` is intended to be embedded within a webpage via an `<iframe>`:

```html
<iframe
  src="https://londongophers.github.io/raffle/"
  style="border:0px;width:100%;overflow:hidden">
</iframe>
```

The content of the generated tweet is controlled by two URL parameters:

* `greeting` - the text to show before the hash
* `hashtags` - comma-separated list of Twitter hashtags that should follow the hash

Both the greeting and hashtags should be URL-encoded. For example, the London Gophers URL is:

https://londongophers.github.io/raffle/?greeting=Hey%20%40LondonGophers%2C%20please%20enter%20me%20into%20the%20%40jetbrains%20raffle!&amp;hashtags=LondonGophers

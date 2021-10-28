# PKSUID - [KSUID](https://github.com/segmentio/ksuid) with arbitrary string prefix

PKSUID is a small extension of the KSUID (K-Sortable Globally Unique IDs) 
[https://github.com/segmentio/ksuid](https://github.com/segmentio/ksuid) which allows prefixing KSUID with the
16 bytes arbitrary string to achieve the following kind of the identifiers:
```text
user:208AocW9380pRVFbQ6sj6Oofu0w
product:208BJB4ntOzehBIADBW6h69juir
anything208BLsTohQ0OnynbQxoeyoRu086
```

Stripe uses this kind of identifier.

## Why using the prefix?

The prefix part might carry additional information valuable for an application.
For example the one can indicate the polymorphic relationship of a resource using an identifier with a prefix.

If you don't need prefixes, we strongly recommend using [KSUID](https://github.com/segmentio/ksuid)

## Example

```go
	userPrefix := pksuid.Prefix{'u', 's', 'e', 'r', ':'}
	id := pksuid.New(userPrefix)

	fmt.Println(id)
	// prints something like this: user:208ETB1zVhm50luzEgNAnLRNQZE

	id2, _ := pksuid.Parse("user:208ETB1zVhm50luzEgNAnLRNQZE")
	fmt.Println(id2.Prefix())
	fmt.Println(id2.Prefix() == userPrefix)
	fmt.Println(id2.ID())
	// user:208Np9GX0sC0J42peUKj4nW1CW9
	// user:
	// true
	// 208ETB1zVhm50luzEgNAnLRNQZE
```


# pyraset

[![Build Status](https://travis-ci.com/gofunky/pyraset.svg?branch=master)](https://travis-ci.com/gofunky/pyraset)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/0b1c4be278ab4d11a8f9658727cdfbb1)](https://www.codacy.com/app/gofunky/pyraset?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=gofunky/pyraset&amp;utm_campaign=Badge_Grade)
[![GoDoc](https://godoc.org/github.com/gofunky/guidelines?status.svg)](https://godoc.org/github.com/gofunky/pyraset)
[![Go Report Card](https://goreportcard.com/badge/github.com/gofunky/pyraset)](https://goreportcard.com/report/github.com/gofunky/pyraset)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=gofunky/pyraset)](https://dependabot.com)
[![GitHub License](https://img.shields.io/github/license/gofunky/pyraset.svg)](https://github.com/gofunky/pyraset/blob/master/LICENSE)
[![GitHub last commit](https://img.shields.io/github/last-commit/gofunky/pyraset.svg)](https://github.com/gofunky/pyraset/commits/master)

pyraset is an improved version of [golang-set](https://github.com/deckarep/golang-set) that features fast structure and subset support

## How to get it

`pyraset` uses go modules and is currently released under the `v2` path.

```bash
go get -u github.com/gofunky/pyraset/v2
```

## What to regard when migrating from `golang-set`

Before migrating from `golang-set` to `pyraset`, there a few things to consider.
Nevertheless, the necessary code changes make your code more awesome.

1. `pyraset` uses simplified constructors, accepting only variadic arguments, check the [docs](https://godoc.org/github.com/gofunky/pyraset) for details.
2. `Equal` has become deep.
3. In tests, use [cmp](https://github.com/google/go-cmp) instead of `reflect` for comparisons.

## Why another set library

Before `pyraset`, I had been using `golang-set`. It works fine for simple operations on small sets.
However, I noticed that some `Equal` calls did not work as I expected them to work.

Take this example:

```golang
a := NewSet("test")
b := NewSet("test")
subSetA := NewSet("foo")
subSetB := NewSet("foo")

// They should be equal, since both contain "test"
fmt.Printf("are a and b equal: %v\n", a.Equal(b))
// are a and b equal: true -> pass

a.Add(subSetA)

// They should not be equal, only a contains subSetA
fmt.Printf("are a and b equal: %v\n", a.Equal(b))
//  are a and b equal: false -> pass

b.Add(subSetB)

// Shouldn't they know be equal? Both are Set { "test", Set { "foo" } }
fmt.Printf("are a and b equal: %v\n", a.Equal(b))
//  are a and b equal: false -> surprise

// Let's try to use subSetA instead
b.Remove(subSetB)
b.Add(subSetA)

fmt.Printf("are a and b equal: %v\n", a.Equal(b))
//  are a and b equal: true
```

What this example shows, apparently, is that the set doesn't use deep equality checks.
For interface types, that is no surprise, for subsets, however, it is indeed unexpected because even `Power` creates subsets.

In `pyraset`, this is solved by using hashes. So basically, instead of a comparable map as used in `golang-set`, we use a hashmap.
Hashing, however has a performance tradeoff. `CRUD` operations are more complex, set operations, in contrast, are performing better.
Some operations perform multiple magnitudes better than `golang-set`.

Hashing gives the user another benefit. `pyraset` uses `hashstructure` and `xxhash` by default.
Accordingly, one may know compare any struct or interface that `hashstructure` accepts as it was a deep equals.

With `Contains`, one may notice that a comparable map performs much better.
Hashing every given element, has a cost. Other operations that rely on `Contains` also have this bottleneck.
To solve this, `pyraset` uses a comparable cache map that maps the generated hashes.
The overhead for the cache map is fair and, in turn, improves the performance of all operations
that compare externally given elements.

## How to use it

In the progress of refactoring `golang-set`, some new operations have been introduced,
others have been dropped.

There are only three constructors anymore. All three accept variadic arguments, no slices.
These are the constructors:
* `NewSet(...)` for thread-safe and cached sets
* `NewUnsafeSet(...)` for non-thread-safe but cached sets
* `SetOptions.New(...)` for sets that are based on individual options (e.g., disabled caching)

`Set.Add(...)` and `Set.Remove(...)` also accept variadic arguments.

Furthermore, there are new "converters" for thread safety:
* `Set.CoreSet()` to receive the non-thread-safe version from any set
* `Set.ThreadSafe()` to receive the thread-safe version from any set

The calculated hash representation can be accessed via `Set.Hash()`.

### UpdateHash

`pyraset` calculates the hash only once for every element to improve the performance.
However, if someone added a pointer to the set and modified the referred struct, `pyraset` will not notice the changes implicitly.
If the elements aren't immutable, one has to call `Set.UpdateHash()` explicitly if the corresponding set operations should consider possible changes.
The update call is optional. If the set operation results only depend on the original values, there is no update necessary.

#### Example

```golang
a := NewSet("test")
b := NewSet("test")
subSetA := NewSet("foo")
subSetB := NewSet("foo")
a.Add(subSetA)
b.Add(subSetB)

fmt.Printf("are a and b equal: %v\n", a.Equal(b))
// are a and b equal: true -> check

subSetA.Add("bar")
subSetB.Add("bar")
fmt.Printf("are a and b equal: %v\n", a.Equal(b))
// are a and b equal: false -> oh no, what's wrong?

updatedA := a.UpdateHash()
updatedB := b.UpdateHash()

fmt.Printf("are a and b equal: %v\n", a.Equal(b))
// are a and b equal: true -> check
fmt.Printf("number of updates:\na:%v\nb:%v\n", updatedA, updatedB)
// number of updates:
// a: 1
// b: 1
// -> check
```

### Transformed examples as known from `golang-set`

```golang
requiredClasses := pyraset.NewSet()
requiredClasses.Add("Cooking")
requiredClasses.Add("English")
requiredClasses.Add("Math")
requiredClasses.Add("Biology")

scienceClasses := pyraset.NewSet("Biology", "Chemistry")

electiveClasses := pyraset.NewSet()
electiveClasses.Add("Welding")
electiveClasses.Add("Music")
electiveClasses.Add("Automotive")

bonusClasses := pyraset.NewSet()
bonusClasses.Add("Go Programming")
bonusClasses.Add("Python Programming")

// Show me all the available classes one can take
allClasses := requiredClasses.Union(scienceClasses).
    Union(electiveClasses).Union(bonusClasses)
fmt.Println(allClasses)
// Set{Cooking, English, Math, Chemistry, Welding, Biology, Music, Automotive, Go Programming, Python Programming}


// Is cooking considered a science class?
fmt.Println(scienceClasses.Contains("Cooking"))
// false

//Show me all classes that are not science classes.
fmt.Println(allClasses.Difference(scienceClasses))
// Set{Music, Automotive, Go Programming, Python Programming, Cooking, English, Math, Welding}

// Which science classes are also required classes?
fmt.Println(scienceClasses.Intersect(requiredClasses))
// Set{Biology}

// How many bonus classes do you offer?
fmt.Println(bonusClasses.Cardinality())
// 2

// Do you have the following classes? Welding, Automotive and English?
fmt.Println(allClasses.IsSuperset(pyraset.NewSet("Welding", "Automotive", "English")))
// true
```

## Performance

The following charts show the ns/op ratios for the different set operations contrasted to `golang-set`.
The smaller the values, the better the scores are, proportionally to the other ones.

![safety](/docs/assets/img/pyraset_safety_performance.png)

![safety](/docs/assets/img/pyraset_ns_op.png)

More benchmark results can be accessed any analyzed [here](https://plot.ly/~mafax/116/).

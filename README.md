#Baduk Readme

Yo dawg, I heard you liked Go, so I put Go in your Go so you could Go on the Go.

##What?

I like Go (the game, also known as Baduk) and Go (the programming language) and decided to build this support library so I could build a Go-webapp to play Go (hence the Yo-dawging). 

But mainly, I built this with the express (and eventual) purpose to play Go with my former (and future?) roommate, and he's much more amenable to my ideas when they involve wordplay, puns, recursive thinking, and self-references. (hi Steve!)

##Usage

Import the package like any remote repo in Go:

```go
import "github.com/acityinohio/baduk"
```

Then initialize a baduk.Board

```go
var b baduk.Board
err := b.Init(13) //Size can be anywhere from 3 to 19, inclusive
//don't forget your errors!
if err != nil {
	//deal with err
}
```

You can set pieces using (x,y) coordinates, anywhere from (0,0) to (size-1, size-1)

```go
var err error
err = b.SetB(0,0)   //Sets black piece at 0,0
err = b.SetW(12,12) //Sets white piece at 12,12
if err != nil {
	//dealWithIt.gif
}
```

My favorite part (and what will help eventually with the whole web app shtick) is that every state of the board can be Encoded/Decoded into a URL-friendly base64-encoded string. Check it:

```go
var err error
err = b.Encode() //Encodes base64-string into b.State
if err != nil {
	//dealWithIt.gif
}
fmt.Println(b.State) //Now do whatever you want with the string
var c baduk.Board
err = c.Decode(b.State) //...like create a new board with the same state
```

This will be super useful for ephemeral, webapp-based games with canonical URLs.

##Testing

You can test the package with the handy "go test" command---the tests are included in the baduk_tests.go package.

##To Do:

A whole lot. Scoring isn't working yet, which is kind of necessary. More test coverage is needed. Just a whole bunch of everything.

##Why not call this package Go?

Yes, we all get off on Xzibit memes, but that goes too far, even for me. The name collisions alone would drive me bonkers.

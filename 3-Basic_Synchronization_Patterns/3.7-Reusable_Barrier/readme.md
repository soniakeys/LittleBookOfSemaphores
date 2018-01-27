# 3.7 Reusable Barrier

Oh, this section seems a mess to me.  The situation first presented is that
of a barrier inside a loop, but then the whole loop thing is never mentioned
again.  Then the puzzle is about this turnstile that I earlier dismissed as a
contrivance.

## Book solution

There's so much that makes me squirm here I'm not even going to implement it.
The text talks about a two-phase barrier, but haven't we been needing to do
this all along?  We needed in a number of programs so far to use a WaitGroup
to wait for all goroutines to complete before exiting main, in addition to
whatever other synchronization problem we were solving.

The solution code shows a line

    \# rendezvous

and a line

    \# critical point

But this is no more than was in the Barrier example.  What are we showing
that is new?  Apparently this trick of counting up and then counting down.

Enough.

## Reformulation

Let's reformulate this problem and try to do something more sensible.
I like the initially mentioned challenge of putting a barrier in a loop.
Let's do that, and simplify by having each worker output a message, then
reach the barrier.

## mainLoop.go

WaitGroups are reusable.  This was used in some previous examples.  The
example here is almost dead simple.  A worker just does some work -- printing
a message, and waits on the waitgroup.  Main runs a loop for any number of
cycles where it starts workers and waits for them to finish.  WaitGroups are
a kind of barrier.  They're pretty easy to reason about and use.

## barrier.go

In mainLoop.go, we dispensed with the channel close idiom even.  Can we revisit
the original barrier problem and simplify it in the same way?  Barrier.go is
it.  Just two kinds of work, each using the WaitGroup to signal completion.

It's starting to feel trivialized though.  What was the big deal in the last
chapter with the barrier in the worker goroutine?  The missing point is that
the barrier technique would be needed if the worker goroutines needed any kind
of continuity across the different kinds of work -- if they accumulated some
kind of per-worker result for example.  Unfortunately that was not given as a
requirement in the last section.  We jumped through some hoops without really
having a reason.

## Iterate on reformulation

Let's try again on making this problem really need a reusable barrier.  Say
there are multiple workers, that each must do some work then reach a barrier
then perhaps repeat, and that workers must independently accumulate some result
over their barrier-punctuated work cycles.

Now it's getting much more interesting.  Say main controls the number of cycles
and must communicate to the workers whether to do another cycle or to wrap up
their work.  The idea is to approach a technique that might be useful in
something like a physics simulation, where workers all compute something for
a new physical state staring from an existing physical state.

## workerLoop.go

The latest reformulation sounds interesting, but really we have no real problem
to solve.  workerLoop.go is just another toy program without much purpose but
anyway shows off some programming techniques.  Anyway, a quick description:

The part of the reformulation that says main "must communicate to the workers
whether to do another cycle or to wrap up their work" kind of implies that
workers must be ready to receive and act on two different messages.  The Go
select statement handles this.

Main sends the allDone message if it needs no more computation cycles.  This
is the signal for the workers to wrap up and terminate.  We implement "wrap
up" as printing our accumulated count, this simple count being our result
accumulated over the work cycles.

Main sends the cycleStart message to send the workers off on one more
computation cycle.

Both the allDone and the cycleStart messages are broadcast messages implemented
with the channel close technique.

Finally, the barrier between computation cycles is a little tricky.  Closing
a channel is not a "reusable" operation.  That is, there is no way to reopen
a closed channel.  Channels are super cheap to create though.  No need to
reuse, just re-make the channel.  The last trick is to use one more channel,
cycleReset, as the final barrier of a cycle.

----
2018/01/27 16:46:14 cycle 0
2018/01/27 16:46:14   gr 0 working
2018/01/27 16:46:14   gr 1 working
2018/01/27 16:46:14   gr 2 working
2018/01/27 16:46:14 cycle 1
2018/01/27 16:46:14   gr 2 working
2018/01/27 16:46:14   gr 1 working
2018/01/27 16:46:14   gr 0 working
2018/01/27 16:46:14 cycle 2
2018/01/27 16:46:14   gr 2 working
2018/01/27 16:46:14   gr 0 working
2018/01/27 16:46:14   gr 1 working
2018/01/27 16:46:14 cycle 3
2018/01/27 16:46:14   gr 2 working
2018/01/27 16:46:14   gr 1 working
2018/01/27 16:46:14   gr 0 working
2018/01/27 16:46:14 gr 2 counted to 4
2018/01/27 16:46:14 gr 0 counted to 4
2018/01/27 16:46:14 gr 1 counted to 4
----
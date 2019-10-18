# EPIC MATH GAME

`EPIC MATH GAME` is javascript based "game" that asks you to answer trivial multiplication questions. High scores are tracked by the server.

This is a terrible game in that it's not any fun. However, it's a fun game in that if you have no experience hacking browser-based games it's a trivial one to start with.

So, _cheat on_!

## Usage

To compile:

```
go build -o emg cmd/emg/main.go
```

To run:

```
./emg -listen localhost:8888
```

If the `-listen` parameter isn't supplied it defaults to all interfaces, port `8000`.

`Ctrl-C` will stop the server.

Scores are reset when the server restarts.

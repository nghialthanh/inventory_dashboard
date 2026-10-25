# Stock-out, proposed action, and image are separate commands

PATCH corrects only the fields the current stock state allows. Departure, proposed action, and the primary image each have their own route, so each rule lives with one command.

## Considered options

- One PATCH that accepts stock state, dates, action, and image together. The client can send a mix the rules forbid, and the handler has to pick which field failed.

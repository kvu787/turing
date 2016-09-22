usage: turing [-h | -t]

Turing accepts a definition of a Turing machine from standard input and 
simulates its execution.

The -t flag accepts an extended definition that includes a tape and 
initial position for the head.

The format of the definition follows:

  head position    (if -t is specified)
  tape             (if -t is specified)
  blank character
  initial state
  final state
  transition rules

The format of a transition rule is:

  current_state input next_state output move

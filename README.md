# ExamDisys

To run the program you have to open up three different terminals

In the first two terminals, you have to start the two respective servers. They have to be started at port 5001 and port 5002.

Feel free to copy the following three lines to each respective terminal to start the servers:

    go run server/server.go -port 5001
    go run server/server.go -port 5002

In the last terminal you have to open the client on a port of your choosing
To do so feel free to copy the following line:

    go run client/client.go -cPort 4041

It is important to wait for everything to be connected.

In the client's terminal you can add a word and a definition to the dictionary by first writing

    add

to the terminal (followed by 'enter'). Then you can write word you want to add (followed by 'enter') and
lastly the definition to the word (followed by 'enter')

In the client's terminal you can read a definition from the dictionary you write

    read

in the terminal (followed by 'enter') and then the word you want to read (followed by 'enter')

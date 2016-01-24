File Modifications
==================

File modifications are simply modifications to be performed on the file data stream as it is read from disk, before it is handed to the backup engines.

## Interface

The interface includes an Encode and Decode function that takes an io.Reader and returns an io.Reader. Encode obviously performs the requested action (e.g. compress or encrypt), and Decode reverses it. Modifications are called in a specific order for backing up files, which is then reversed when restoring them.

#### Definition

Modifications are defined by their name and options, and must implement a Configure function which accepts the the provided options and configures itself, or returns an error

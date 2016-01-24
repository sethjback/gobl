Engines
==========

Engines are responsible for handling the actual saving/retrieval of the data, though that only really entails implementing the backup and/or restore interfaces (for example the logger engine only implements backup and simply counts the bytes in the stream then makes a log entry about the file)

## Interface

#### Backup

The backup interface requires that an engine be able to backup and retrieve a file, as well as confirm that a given signature needs to be backed up (i.e. the engine doesn't already have a copy of that file saved).

#### Restore

The restore interface basically just needs to take a file reader and save the bytes piped through it

## Definition

Engines are defined by a struct that contains their name, and a map of options.

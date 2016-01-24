# gobl

gobl is a client/server host-based backup solution implemented in go. The intent behind it is to be simple in the sense that it doesn't implement things like disk-safes, continuous data protection, or worry about tape based media. The express goal is to be able to run scheduled differential backups where the file is read modified and saved in a single pass.

The project mainly came about out of my desire to explore the go language, and the need for a backup solution that would fit within the AWS context in which I would use it.

For the general idea behind the architecture see the docs section.

#### Disclaimer

This is **alpha** software. It is **NOT** finished, and it is an initial attempt at coding in go so there are no guarantees that the software is idiomatic, bug free, efficient, etc. Think of it more as the thumbnail sketch of where I would like the project to go, but it is a long way from here to complete.

#### Roadmap

The high-priority items are:

* test coverage (non-existent atm!)
* consistent logging
* better error handling
* api definitions with swagger
* job expiration with auto removal
* agent removal
* file browser implementation
* documentation update

Next-up features/improvements

* compression engine intelligently ignoring already compressed file formats
* S3 backup engine
* Web front end
* encryption modification

At some point in the future

* saving file metadata (owner/group/permissions)
* user authentication
* remote file engine (would allow for centralized backup server)
* easy install packages
* additional database modules (e.g. mysql)

## License

Copyright 2016 Seth back

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

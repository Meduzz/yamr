# yamr
Yet another maven repo.

The end game here would be a system that can be ran in docker and scaled in any direction.

## Variables

Bits of this software are controlled by ENV variables. Others by flags you pass when you start the app.

### Select storage engine.

By setting the flag ```storage``` you control the which storage engine to use.
At the moment, only one storage engine are available, and that's ```filesystem``` that is also the default value for this variable.

#### Filesystem module settings

The filesystem storage engine also accepts an ENV variable called ```FS_PATH``` that controls where to store the files. The default value are ```files```.
An attempt will be made to create this dir if it does not already exists.

### Metadata

Metadata about files are stored in a pgsql, this is a feature that are not optional. 
This is so that the actual storage engine dont have to worry about doing unnecessary (and perhaps expensive) lookups. 
Under the hood libpg are used, and all ENV variables it can handle can (and should) be used to connect to your pg-database.
[Postgres ENV connection reference](https://www.postgresql.org/docs/9.6/static/libpq-envars.html)
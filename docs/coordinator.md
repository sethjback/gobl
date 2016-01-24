Coordinator
==========

The Coordinator is designed to be the central point of management for the backup process. It handles all the configuration for the backup jobs and maintains all the records of their execution.

The Coordinator handles 4 primary things:
1. Agents
2. Backup jobs
3. Job Records
4. File Catalog

## Agents

In order to configure a backup job the Coordinator must first know about the Agent. Agents are identified by their public key, which is used to sign every request made to the Coordinator's API. Adding an agent is really the process of creating an entry in the database along with the Agent's key.

## Backup Jobs

Backup jobs are descriptions of backups to be made.

## Job Records

## File Records

File records are a way for the Agent to know if a file has already been backed up, or if it needs to make a fresh copy of it. Each record is stored under the key of the agent that creates it, and is a base64URL hash of the following JSON data:

```json
{
  "name": "filename",
  "path": "full/file/path",
  "hash": "imohash of the file",
  "modifications": "modifications"
}
```
Using this information, it is possible in a single query to know if a given file has already exists somewhere in the backup file set.

It is important to note that the Coordinator cannot guarantee that the backup engine has indeed saved the file correctly, where the engine has saved the file, or that the engine has not subsequently removed the file. All of that logic is the domain of a well behaved backup engine, and the Coordinator only keeps track of what the backup engines tell it.

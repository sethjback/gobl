Architecture
============

There are two main components to Gobble:
* Coordinator
* Agents

### Coordinator

The Coordinator is brains of the operation. It maintains a list of active Agents, the backup schedules and status, file lists, and logs. The Coordinator implements a restful API that is used both to configure the Coordinator and by the Agents to update the Coordinator on job status.

### Agents

Agents run on each server to be backed up and are directly responsible for handling the entire backup process of that server. The agents implement a restful API that the Coordinator communicates with to initiate a job (backup or restore). In turn, the Agents contact the Coordinator's API to update it on the job status. All the work of backing up the actual files (whatever that may entail) is handled by the Agent.

### Workflow

All configuration is done on the the Coordinator. Backups are handled as Jobs that contain the following:

* Target Agent: the agent to run the job on
* Files to be backed up
* In-line file modification(s)
* Backup engine(s) to use

Gobl was intended to be flexible: there can be any number of in-line file modifications which are called sequentially and can modify the data after it is read from disk but before it is written to the backup engine. Likewise, there can be any number of backup engines, each with their own configuration, and the modified file stream will be passed to each.

When the job runs (according to the schedule), the Coordinator connects to the specified Agent and hands it the job's specs. The Agent then proceeds with the backing up the file using the specified modifications and engines. When the file has been successfully saved, the Agent notifies the Coordinator, which then adds it to the file records for that backup.


### File Signatures

Before each file is backed up a signature is created which is then used to uniquely identify that file in it's current state. The signature consists of the following:

* filename
* path
* file hash
* modifications

Before a file is backed up, this signature is sent to the engine, which is then responsible to see if it already has a file with this exact signature saved. If so, that engine is skipped (though the file is still included in the job record). This is essentially the differential part of the backup process.

The file hash is a combination of the file size and samples from the beginning, middle, and end of the file. This saves the need to hash an entire file if it is large (though small files are completely hashed).


### Job Records

A job record is stored on the Coordinator and keeps track of the actual files (signatures) backed up or restored by the agent. Restores are done from backup job records: the job record contains a list of the modifications and engines used in the backup, which is then used to retrieve the backed up file and reverse the modifications.

### Communication

The Coordinator and Agents communicate via calls to their respective RESTful apis.  To ensure the request is authentic and hasn't been modified, each sensitive action (for example starting jobs/restores, adding files to jobs, updating job status) is signed an authenticated using RSA PSS.

Gobl doesn't generate the keys to use, but expects to be configured with a key path for both the Coordinator and Agent. To start the process, the Coordinator's public key must be manually installed on the Agent. When an Agent is added to the Coordinator, it queries the Agent and asks for it's public key.

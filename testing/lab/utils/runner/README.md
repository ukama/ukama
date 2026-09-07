# Distributed P0 runner

This directory implements a deliberately small distributed test runner:

- the local machine is the controller;
- S3 stores input and output files;
- each EC2 instance receives one text file of scenario paths;
- each worker runs its scenarios sequentially;
- workers upload status and results, then terminate;
- the local controller displays status and merges the reports.

There is no queue, scheduler, database, daemon, web service, Lambda, Batch,
or cross-worker communication.

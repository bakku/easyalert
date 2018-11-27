# Features

Easyalert should behave according to its name. Sending alerts in an easy and fast way without needing any libraries that have to be used or any executables that need to be installed. It should expose an alert mechanism via a simple REST-like HTTP API.

- Each step a user can take should be achievable via the REST API and a web interface
- The REST API should provide helpful information for executing other API requests via command line
- Users can create accounts and edit them
- Users can send emails to their email address
- Users can view their last emails and see whether their delivery was successful
- Email body is confidential and won't be stored inside the database and won't be part of logging

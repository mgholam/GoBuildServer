# Go Build Server

A simple web server to build any project via scripts for Continuous Integration/Deployment systems.

## Config file

 ```
 {
   "Port": 5000,
   "Projects": [
     {
       "Name": "Project 1",
       "Status": "",
       "LastBuildDate": "2022-02-17T09:45:17",
       "CmdPath": "d:/project/build.cmd",
       "ErrorPath": "d:/project/errors.txt",
       "LastBuildDuration": 0
     }
   ]
 }
 ```
- **Port** : the servers web port
- **CmdPath** : the platform specific build script to build the project and create `errors.txt` file on errors
- **ErrorPath** : the path to the `errors.txt` file for build errors



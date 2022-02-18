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
       "ErorrPath": "d:/project/errors.txt",
       "LastBuildDuration": 0
     }
   ]
 }
 ```




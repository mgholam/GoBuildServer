# Go Build Server

A simple web server to build any project via scripts for Continuous Integration/Deployment systems.

## Config file

 ```json
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

### Example build script for .net
The following `build.cmd` file will build 3 solutions in order and stop if any fails outputing the errors to `errors.txt`

```
msbuild  project1\All.sln -v:m -flp1:logfile=errors.txt;errorsonly
if errorlevel 1 goto errorDone

msbuild  project2\All.sln -v:m -flp1:logfile=errors.txt;errorsonly
if errorlevel 1 goto errorDone

msbuild  project3\all.sln -v:m -flp1:logfile=errors.txt;errorsonly
if errorlevel 1 goto errorDone


:errorDone

```



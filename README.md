# SmartRecipes
This is tha analytical tool which reads the file and provides the analytical results.

## Pre-requites  
Go 1.16  
VSCode 1.56  
Go plugins  

## Build this tool  
Navigate to root folder  
RUN "go mod download" --> to download dependencies  
RUN "go build -o main ." --> to build the code and create binary  

## Run this tool  
Since this tool uses some custom inputs for queries  
Below Env variables are used to custom the queries  
MATCH_BY_NAME="<queries recipes based on names separated by space>" eg ["Apple Pizza"]  
QUERY_POSTCODE="<queries delivery count based on postcode>" eg "10001"  
QUERY_DELIVERY_TIME="<queries delivery count based on time>" eg "7AM - 10PM"  
FILE="<filename>" eg test_data.json  

You can run this tool with or without docker  

## Without Docker  
1. First set these env variables described above  
2. RUN "./main  <location to filename>" eg "~/data.json" [env variable FILE can also be used to set file to run]  

## With Docker  
The json files we want to run has to be in working directory  
else it would not be able to access that from container  
RUN "docker build -t recipe-app ." --> this would create image of this cli  
RUN docker run --env FILE="<filename>" --env QUERY_DELIVERY_TIME="<query delivery name in format '3AM - 12PM'>" --env QUERY_POSTCODE="<query postcode>" --env MATCH_BY_NAME="query to match by name separated by space" recipe-app  

Eg for passing env variables  
docker run --env FILE="test_data.json" --env QUERY_DELIVERY_TIME="12AM - 10PM" --env QUERY_POSTCODE="10001" --env MATCH_BY_NAME="Apple" recipe-app  

## Note 
1. QUERY_POSTCODE and QUERY_DELIVERY_TIME need to be used together to query delivery count per postcode and time  
2. If using file name as command args, then env variable FILE wont be used as command would get more precedence over env  
3. By default MATCH_BY_NAME = "Potato Veggie Mushroom"  
4. By default QUERY_POSTCODE="10120" and QUERY_DELIVERY_TIME="10AM - 3PM"  

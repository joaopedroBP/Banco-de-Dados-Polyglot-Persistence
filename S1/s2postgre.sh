#!/bin/bash

set -e

cd ..
cd S2_Postgre
javac -cp "postgresql-42.7.3.jar" Main.java
java -cp ".:postgresql-42.7.3.jar" Main $1 $2 $3 $4 $5 $6 $7 $8 $9
cd ..
cd S1

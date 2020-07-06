# SensorReadWebsite

Website contains database <--> backend <--> frontend. It represents a website where operetor can enter,update,get,delete sensor measured data.
![107077573_300673567642198_3677948170768276020_n](https://user-images.githubusercontent.com/62447953/86564483-9aa07d80-bf66-11ea-9050-0cf5c7fceadf.png)
## Cloning
```
git clone https://github.com/aTTiny73/SensorReadWebsite.git
```
## Database setup

To setup database first you need to instal mysql-server.
To get the exact same table as me, inside the mysql shell, type these commands :
```
CREATE DATABASE SENSORDATA;
USE SENSORDATA;
CREATE TABLE READINGS
(
    ID varchar(255) NOT NULL,
    Temperature varchar(255) NOT NULL,
    Humidity varchar(255) NOT NULL,
    CO2 varchar(255) NOT NULL,
    Time varchar(255) NOT NULL,
    PRIMARY KEY (ID)
);
```
Now its time to setup user you can do that by running this command in mysql shell:

```
CREATE USER 'testuser'@'localhost' IDENTIFIED BY 'testuser';
```
Now we need to grant all privileges to user so he can add to the tabel delete etc.
```
GRANT ALL PRIVILEGES ON SENSORDATA.READINGS TO 'testuser'@'localhost';
```

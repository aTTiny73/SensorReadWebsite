# SensorReadWebsite

Website contains database <--> backend <--> frontend. It represents a website where operetor can enter,update,get,delete sensor measured data.
![Screenshot from 2020-07-07 11-43-56](https://user-images.githubusercontent.com/62447953/86762900-3dc8c400-c047-11ea-87e6-5cb144202154.png)
## Cloning
```
git clone https://github.com/aTTiny73/SensorReadWebsite.git
```
## Database setup

To setup database first you need to install mysql-server.
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
CREATE USER 'testuser'@'localhost' IDENTIFIED BY 'testpassword';
```
Now we need to grant all privileges to user so he can add to the tabel delete etc.
```
GRANT ALL PRIVILEGES ON SENSORDATA.READINGS TO 'testuser'@'localhost';
```

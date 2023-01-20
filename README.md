# POC InfluxDB
Im going to try influxdb on some sensors data.

#### Loading Data
- download data from [pump-sensor-data kaggle](https://www.kaggle.com/datasets/nphantawee/pump-sensor-data) into `resources/`
  ![](resources/dataset.png)
- seperate datasets files for each sensors and machine status as I wish to simulate multiple measurements.
- create measurements for all sensors as format (machine,sensor=sensor1 value=xxxx timestamp)
- Load data via channel and goroutine.

#### Install InfluxDB
following download guide in `https://portal.influxdata.com/downloads/`

#### Integrate with InfluxDB
- create measurements.
- stream data to measurement.
- create query all data with mean values in minute.

#### Definition of Measurement on platform
- Define Measurement name
- Define tags
- Upload file with fields as columns and value in rows

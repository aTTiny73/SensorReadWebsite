function loadAllData() {
    var table = document.getElementById("dataTable");
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
      if (this.readyState == 4 && this.status == 200) {
        var data = JSON.parse(this.responseText);
        console.log(data)
        for(var i = table.rows.length - 1; i > 0; i--)
        {
            table.deleteRow(i);
        }
        for(var i = 0;i < data.length;i++){
            
            var row = table.insertRow(-1);
            var cell1 = row.insertCell(0);
            var cell2 = row.insertCell(1);
            var cell3 = row.insertCell(2);
            var cell4 = row.insertCell(3);
            var cell5 = row.insertCell(4);
            
           cell1.innerHTML = data[i].id;
           cell2.innerHTML = data[i].temperature;
           cell3.innerHTML = data[i].humidity;
           cell4.innerHTML = data[i].co2;
           cell5.innerHTML = data[i].time; 
         }
      }
    };
    xhttp.open("GET", "http://localhost:8090/getReadings", true);
    xhttp.send();
  }
  
  function loadData() {
    var str = document.getElementById("input1").value;
    var table = document.getElementById("dataTable");
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function() {
      if (this.readyState == 4 && this.status == 200) {
       var data = JSON.parse(this.responseText);
       console.log(data);
       for(var i = table.rows.length - 1; i > 0; i--)
        {
            table.deleteRow(i);
        }
       var row = table.insertRow(1);
       var cell1 = row.insertCell(0);
       var cell2 = row.insertCell(1);
       var cell3 = row.insertCell(2);
       var cell4 = row.insertCell(3);
       var cell5 = row.insertCell(4);
  
       cell1.innerHTML = data[0].id;
       cell2.innerHTML = data[0].temperature;
       cell3.innerHTML = data[0].humidity;
       cell4.innerHTML = data[0].co2;
       cell5.innerHTML = data[0].time;
        
      }
    };
    xhttp.open("GET", "http://localhost:8090/getReading/"+str, true);
    xhttp.send();
  }
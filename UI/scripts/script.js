// Select container section to appear when clicked
var telemetryInputContainer = document.querySelector('.telemetry-input-container');
var coordinateInputContainer = document.querySelector('.coordinate-input-container');
var angleInputContainer = document.querySelector('.angle-input-container');

// Set buttons to be able to do actions
var setTelemetryBtn = document.getElementById('setTelemetryBtn');
var setCoordinateBtn = document.getElementById('setCoordinateBtn');
var setAngleBtn = document.getElementById('setAngleBtn');

var sendCommandBtn = document.getElementById('sendCommandBtn');
var sendCommandBtn2 = document.getElementById('sendCommandBtn2');
var sendCommandBtn3 = document.getElementById('sendCommandBtn3');

var getTelemetryBtn = document.getElementById('getTelemetryBtn');

var script1Btn = document.getElementById('script1Btn');

// Don't display the containers by default
telemetryInputContainer.style.display = 'none';
coordinateInputContainer.style.display = 'none';
angleInputContainer.style.display = 'none';

// Add click event listener to the Change Telemetry button
setTelemetryBtn.addEventListener('click', function() {
    // Toggle the display property of the telemetry input container
    if (telemetryInputContainer.style.display === 'none' || telemetryInputContainer.style.display === '') {
        angleInputContainer.style.display = 'none';
        coordinateInputContainer.style.display = 'none';
        telemetryInputContainer.style.display = 'block';
    } else {
        telemetryInputContainer.style.display = 'none';
    }
});

// Add click event listener to the Change Coordinate button
setCoordinateBtn.addEventListener('click', function() {
    // Toggle the display property of the coordinate input container
    if (coordinateInputContainer.style.display === 'none' || coordinateInputContainer.style.display === '') {
        angleInputContainer.style.display = 'none';
        telemetryInputContainer.style.display = 'none';
        coordinateInputContainer.style.display = 'block';
    } else {
        coordinateInputContainer.style.display = 'none';
    }
});

// Add click event listener to the Change Angle button
setAngleBtn.addEventListener('click', function() {
    // Toggle the display property of the angle input container
    if (angleInputContainer.style.display === 'none' || angleInputContainer.style.display === '') {
        telemetryInputContainer.style.display = 'none';
        coordinateInputContainer.style.display = 'none';
        angleInputContainer.style.display = 'block';
    } else {
        angleInputContainer.style.display = 'none';
    }
});

// Event listener for requesting telemetry from Space CNDH
getTelemetryBtn.addEventListener('click', function() {

    statusIsGood(function (isGood) {
        if (isGood) 
        {
            const request = new XMLHttpRequest();
            const url = 'http://localhost:8080/requestTelemetry'
            request.open("GET", url);
            request.send();
            // Send alert if request is successful
            request.onreadystatechange = function() {
                if(this.readyState == 4 && this.status == 200) {
                    window.alert('Request for telemetry has been sent.');
                }
            }
        }
        else 
        {
            window.alert('Status of the link was not open');
        }
    });
    
});

var update = new EventSource('/update');
update.onmessage = event => {
    console.log('Updated telemetry');
    const telemetry = JSON.parse(event.data);
    console.log('Received JSON: ', telemetry);
    document.getElementById('x-coordinate').innerHTML = telemetry.coordinate.x + ', ';
    document.getElementById('y-coordinate').innerHTML = telemetry.coordinate.y+ ', ';
    document.getElementById('z-coordinate').innerHTML = telemetry.coordinate.z;
    document.getElementById('pitch').innerHTML = telemetry.rotation.p + '°, ';
    document.getElementById('yaw').innerHTML = telemetry.rotation.y + '°, ';
    document.getElementById('roll').innerHTML = telemetry.rotation.r+ '°';
    document.getElementById('temp').innerHTML = telemetry.temp + ' °C';
    document.getElementById('payload-power').innerHTML = 'Payload: ' + telemetry.status.payloadPower;
    document.getElementById('data-waiting').innerHTML = 'Waiting for Data: ' + telemetry.status.dataWaiting;
    document.getElementById('charge-status').innerHTML = 'Charge Status: ' + telemetry.status.chargeStatus;
    document.getElementById('voltage').innerHTML = 'Current Voltage: ' + telemetry.status.voltage;
};

window.addEventListener('beforeunload', function() {
    update.close();
});


// Event listener for setting telemetry to Space CNDH
sendCommandBtn.addEventListener('click', function() {

    var X = document.getElementById('coordX').value;
    var Y = document.getElementById('coordY').value;
    var Z = document.getElementById('coordZ').value;

    var pitch = document.getElementById('anglePitch').value;
    var yaw = document.getElementById('angleYaw').value;
    var roll = document.getElementById('angleRoll').value;

    // Create an object with the form data
    var formData = {
        coordinate: {
            x: Number(X),
            y: Number(Y),
            z: Number(Z)
        },
        rotation: {
            p: Number(pitch),
            y: Number(yaw),
            r: Number(roll)
        }
    };

    if (checkInvalidInput(X, Y, Z, pitch, yaw, roll))
    {
        window.alert('All input fields must be between -180 and 360.')
    }
    else
    {
        statusIsGood(function (isGood) {
            if (isGood) {
                sendRequest(formData, 1);
            } else {
                window.alert('Status of the link was not open');
            }
        });
    }
});

// Event listener for requesting coordinates to Space CNDH
sendCommandBtn2.addEventListener('click', function() {

    var X = document.getElementById('coordX2').value;
    var Y = document.getElementById('coordY2').value;
    var Z = document.getElementById('coordZ2').value;

    // Create an object with the form data
    var formData = {
        coordinate: {
            x: Number(X),
            y: Number(Y),
            z: Number(Z)
        },
    };

    if (checkInvalidInput(X, Y, Z))
    {
        window.alert('All input fields must be between -180 and 360.')
    }
    else
    {
        statusIsGood(function (isGood) {
            if (isGood) {
                sendRequest(formData, 2);
            } else {
                window.alert('Status of the link was not open');
            }
        });
    }
});

// Event listener for requesting rotation to Space CNDH
sendCommandBtn3.addEventListener('click', function() {

    var pitch = document.getElementById('anglePitch2').value;
    var yaw = document.getElementById('angleYaw2').value;
    var roll = document.getElementById('angleRoll2').value;

    // Create an object with the form data
    var formData = {
        rotation: {
            p: Number(pitch),
            y: Number(yaw),
            r: Number(roll)
        }
    };

    if (checkInvalidInput(pitch, yaw, roll))
    {
        window.alert('All input fields must be between -180 and 360.')
    }
    else
    {
        statusIsGood(function (isGood) {
            if (isGood) {
                sendRequest(formData, 3);
            } else {
                window.alert('Status of the link was not open');
            }
        });
    }

    
});

function checkInvalidInput(...values)
{
    for (var i = 0; i < values.length; i++)
    {
        var numValue = Number(values[i]);
        if (numValue < -180 || numValue > 360 || values[i] === "")
        {
            return true;
        }
    }
    return false;
}

function sendRequest(formData, id)
{
    console.log('called');
    const request = new XMLHttpRequest();
    const url = 'http://localhost:8080/settelemetry?id=' + id;
    request.open("PUT", url);
    console.log('send')
    request.send(JSON.stringify(formData));
    
    // Send alert if request is successful
    request.onreadystatechange = function () {
        if (this.readyState == 4) {
            if (this.status == 200) {
                window.alert('Request to set telemetry has been sent.');
            } else {
                window.alert('Failed to send telemetry request. Status: ' + this.status);
            }
        }
    };
}

function statusIsGood(callback)
{
    const request = new XMLHttpRequest();
    const url = 'http://localhost:8080/status'
    request.open("GET", url);
    request.send();

    request.onreadystatechange = function () {
        if (this.readyState == 4) {
            if (this.status == 200) {
                console.log(this.status)
                callback(true);
            } else {
                console.log(this.status)
                callback(false);
            }
        } 
    }
}

// Event listener for Script1 button
script1Btn.addEventListener('click', function() {

    statusIsGood(function (isGood) {
        if (isGood) 
        {
            const request = new XMLHttpRequest();
            const url = 'http://localhost:8080/execute/Script1'
            request.open("GET", url);
            request.send();
     
            // Send alert if request is successful
            request.onreadystatechange = function () {
                if (this.readyState == 4) {
                    if (this.status == 200) {
                        window.alert('Request was sent successfully.');
                    } else {
                        window.alert('Failed to send request. Status: ' + this.status);
                    }
                }   
            };
        } 
        else 
        {
            window.alert('Status of the link was not open');
        }
    });
});

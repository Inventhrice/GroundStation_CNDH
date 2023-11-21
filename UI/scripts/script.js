// Select container section to appear when clicked
var telemetryInputContainer = document.querySelector('.telemetry-input-container');
var coordinateInputContainer = document.querySelector('.coordinate-input-container');
var angleInputContainer = document.querySelector('.angle-input-container');

// Set buttons to be able to do actions
var setTelemetryBtn = document.getElementById('setTelemetryBtn');
var setCoordinateBtn = document.getElementById('setCoordinateBtn');
var setAngleBtn = document.getElementById('setAngleBtn');
var getTelemetryBtn = document.getElementById('getTelemetryBtn');
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
});

var update = new EventSource('/update');
update.onmessage = event => {
    console.log('Updated telemetry');
    const telemetry = JSON.parse(event.data);
    console.log('Received JSON: ', telemetry);
    document.getElementById('x-coordinate').innerHTML = telemetry.coordinates.x + ', ';
    document.getElementById('y-coordinate').innerHTML = telemetry.coordinates.y+ ', ';
    document.getElementById('z-coordinate').innerHTML = telemetry.coordinates.z;
    document.getElementById('pitch').innerHTML = telemetry.rotations.p + '°, ';
    document.getElementById('yaw').innerHTML = telemetry.rotations.y + '°, ';
    document.getElementById('roll').innerHTML = telemetry.rotations.r+ '°';
    document.getElementById('temp').innerHTML = telemetry.temp + ' °C';
    document.getElementById('payload-power').innerHTML = 'Payload: ' + telemetry.status.payloadPower;
    document.getElementById('data-waiting').innerHTML = 'Waiting for Data: ' + telemetry.status.dataWaiting;
    document.getElementById('charge-status').innerHTML = 'Charge Status: ' + telemetry.status.chargeStatus;
    document.getElementById('voltage').innerHTML = 'Current Voltage: ' + telemetry.status.voltage;
};

window.addEventListener('beforeunload', function() {
    update.close();
});

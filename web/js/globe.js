var POS_X = 1800;
var POS_Y = 0;
var POS_Z = 1800;
var BASE_Y = 0;
var WIDTH = 1000;
var HEIGHT = 600;

var FOV = 45;
var NEAR = 1;
var FAR = 4000;

// some global variables and initialization code
// simple basic renderer
var renderer = new THREE.WebGLRenderer();
renderer.setSize(WIDTH,HEIGHT);
renderer.setClearColorHex(0x000000);

// add it to the target element
var mapDiv = document.getElementById("globe");
mapDiv.appendChild(renderer.domElement);

// setup a camera that points to the center
var camera = new THREE.PerspectiveCamera(FOV,WIDTH/HEIGHT,NEAR,FAR);
camera.position.set(POS_X,POS_Y, POS_Z);
camera.lookAt(new THREE.Vector3(0,0,0));

// create a basic scene and add the camera
var scene = new THREE.Scene();
scene.add(camera);

var light = new THREE.DirectionalLight(0x5555ff, 3.5, 500);

// we wait until the document is loaded before loading the
// density data.
$(document).ready(function()  {
	$('#globe').mousemove(mouseCapture);
	jQuery.get('/distribution.js', function(data) {
		addLights();
		addEarth();
		addClouds();
		addDistribution(data);
		render();
	});
});

function addEarth() {
	var spGeo = new THREE.SphereGeometry(600,50,50);
	var planetTexture = THREE.ImageUtils.loadTexture("/images/world-big-2-grey.jpg");
	var mat2 = new THREE.MeshPhongMaterial({
		map: planetTexture,
		shininess: 0.2
	});
	var sp = new THREE.Mesh(spGeo,mat2);
	scene.add(sp);
}

function addLights() {
	scene.add(light);
	light.position.set(POS_X,POS_Y,POS_Z);
}

function latLongToVector3(lat, lon, radius, heigth) {
	var phi = (lat)*Math.PI/180;
	var theta = (lon-180)*Math.PI/180;
 
	var x = -(radius+heigth) * Math.cos(phi) * Math.cos(theta);
	var y = (radius+heigth) * Math.sin(phi);
	var z = (radius+heigth) * Math.cos(phi) * Math.sin(theta);
 
	return new THREE.Vector3(x,y,z);
}

function addDistribution(data) {
	var geom = new THREE.Geometry();
	var cubeMat = new THREE.MeshLambertMaterial({color: 0x000000, opacity: 0.7, emissive: 0xeeeeff });
	for (var i = 0; i < data.length; i++) {
		var entry = data[i];
		var px = entry.Longitude;
		var py = entry.Latitude;
		var value = entry.ZipCodes / 12;
		var position = latLongToVector3(py,px,600,3);
		var cube = new THREE.Mesh(new THREE.CubeGeometry(5,5,1+value,1,1,1,cubeMat));
		cube.position = position;
		cube.lookAt(new THREE.Vector3(0,0,0));
		THREE.GeometryUtils.merge(geom, cube);
	}
	var total = new THREE.Mesh(geom, cubeMat);
	scene.add(total);
}

var meshClouds = null;
function addClouds() {
	var spGeo = new THREE.SphereGeometry(600,50,50);
	var cloudsTexture = THREE.ImageUtils.loadTexture( "/images/earth_clouds_1024.png" );
	var materialClouds = new THREE.MeshPhongMaterial( { color: 0xffffff, map: cloudsTexture, transparent:true, opacity: 0.1 } );

	meshClouds = new THREE.Mesh( spGeo, materialClouds );
	meshClouds.scale.set( 1.025, 1.025, 1.025 );
	scene.add( meshClouds );
}

function copy(obj) {
	var c = {};
	for (var k in obj) {
		c[k] = obj[k];
	}
	return c;
}

var mouseY = 0;
var mouseX = 0;
function mouseCapture(evt) {
	var y = evt.offsetY;
	mouseY = -(y - (HEIGHT / 2)) * 1.7;
	mouseX = (evt.offsetX - (WIDTH / 2));
}

var lastShiftTime = null;
function render() {
	var timer = Date.now();
	camera.position.x = (Math.cos(timer / 10000) * 1800);
	camera.position.z = (Math.sin(timer / 10000) * 1800);
	if (lastShiftTime === null) {
		lastShiftTime = timer;
	} else {
		var since = (timer - lastShiftTime) * 2;
		var movement = (Math.abs(POS_Y - mouseY) * (POS_Y < mouseY ? 1 : -1)) / since;
		
		POS_Y += movement;
		if (POS_Y < -500) { POS_Y = -500; }
		if (POS_Y > 500) { POS_Y = 500; }
		camera.position.y = POS_Y;
		lastShiftTime = Date.now();
	}
	camera.lookAt(scene.position);
	meshClouds.rotation.y = -(timer/40000);
	light.position = copy(camera.position);
	light.position.y = camera.position.y - 300;
	light.lookAt(scene.position);
	renderer.render(scene, camera);
	requestAnimationFrame(render);
}

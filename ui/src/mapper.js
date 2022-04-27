var map = new ol.Map({
  target: 'map',
  layers: [
    new ol.layer.Tile({
      source: new ol.source.OSM({
        url: '/tile/{z}/{x}/{y}'
      })
    })
  ],
  view: new ol.View({
    center: ol.proj.fromLonLat([37.41, 8.82]),
    zoom: 3,
    minZoom:2,
    maxZoom: 5
  })
});

var stationSource = new ol.source.Vector({
    features: []
  })
var stationLayer = new ol.layer.Vector({
  source: stationSource
});
map.addLayer(stationLayer);


function refreshStations(){
  fetch('/stations')
    .then(response => {
      return response.json();
    })
    .then(stations => {
      stationSource.clear()
      for (station in stations) {
        if (stations[station]['maidenhead'] != "") {
          stationSource.addFeature(new ol.Feature({
            geometry: new ol.geom.Point(ol.proj.fromLonLat([stations[station]['longitude'], stations[station]['latitude']]))
          }))
          console.log(stations[station])
        }
      }
    })
}

refreshStations()
setInterval(refreshStations, 15000)

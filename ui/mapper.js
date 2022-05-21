import Map from 'ol/Map'
import View from 'ol/View'
import Feature from 'ol/Feature'
import Tile from 'ol/layer/Tile'
import OSM from 'ol/source/OSM'
import {Vector as SourceVector}  from 'ol/source'
import {Vector as LayerVector} from 'ol/layer'
import {fromLonLat} from 'ol/proj'
import {Style,Stroke} from 'ol/style'
import {Point,LineString} from 'ol/geom'

var host = ''

if (process.env.NODE_ENV === 'development') { // Or, `process.env.NODE_ENV !== 'production'`
  var host = process.env.HOST
}

var map = new Map({
  target: 'map',
  layers: [
    new Tile({
      source: new OSM({
        url: `${host}/tile/{z}/{x}/{y}`
      })
    })
  ],
  view: new View({
    center: fromLonLat([0.0, 0.0]),
    zoom: 2,
    minZoom: 2,
    maxZoom: 5
  })
});

var stationSource = new SourceVector({
  features: []
})
var stationLayer = new LayerVector({
  source: stationSource
});
map.addLayer(stationLayer);

/* var lineStyle = [new Style({stroke: new Stroke({color: '#d12710', width: 2})})];
 *
 * var rawPointA = fromLonLat([-5.0, 54.0])
 * var rawPointB = fromLonLat([-50.0, 30.0])
 * var pointA = new Feature({geometry: new Point(rawPointA)})
 * var pointB = new Feature({geometry: new Point(rawPointB)})
 *
 * var line = new Feature({geometry: new LineString([rawPointA, rawPointB]),
 *                            style: lineStyle,
 *                            name: "Line :)"})
 *
 * stationSource.addFeature(pointA)
 * stationSource.addFeature(pointB)
 * stationSource.addFeature(line) */

function refreshStations(){
  fetch(`${host}/stations`)
    .then(response => {
      return response.json();
    })
    .then(stations => {
      stationSource.clear()
      for (let station in stations) {
        if (stations[station]['maidenhead'] != "") {
          stationSource.addFeature(new Feature({
            geometry: new Point(fromLonLat([stations[station]['longitude'], stations[station]['latitude']]))
          }))
          console.log(stations[station])
        }
      }
    })
}

refreshStations()
setInterval(refreshStations, 15000)

import Map from 'ol/Map'
import View from 'ol/View'
import Feature from 'ol/Feature'
import Tile from 'ol/layer/Tile'
import OSM from 'ol/source/OSM'
import VectorSource  from 'ol/source/Vector'
import VectorLayer from 'ol/layer/Vector'
import {fromLonLat} from 'ol/proj'
import {Style,Stroke} from 'ol/style'
import {Point,LineString} from 'ol/geom'
import {click} from 'ol/events/condition'
import Select from 'ol/interaction/Select'
import Overlay from 'ol/Overlay'

var host = ''

if (process.env.NODE_ENV === 'development') { // Or, `process.env.NODE_ENV !== 'production'`
  var host = process.env.HOST
}



// ** STATIONS ****************************************
var stationFeatures = {}
var stationSource = new VectorSource({
  features: []
})
var stationLayer = new VectorLayer({
  source: stationSource
});



// ** POPUP ****************************************
const container = document.getElementById('popup');
const content = document.getElementById('popup-content');
const closer = document.getElementById('popup-closer');

const popupOverlay = new Overlay({
  element: container,
  autoPan: {
    animation: {
      duration: 250
    }
  }
})

closer.onclick = function(){
  popupOverlay.setPosition(undefined);
  closer.blur();
  return false;
}


// ** INTERACTION ****************************************
var stationSelectInteraction = new Select({
  condition: click,
  layers: [stationLayer]
})

stationSelectInteraction.on('select', function(e){
  const feature = e.target.getFeatures().item(0);
  const data = feature.values_.data

  content.innerHTML = `${data.callsign}<br>${data.maidenhead}`

  popupOverlay.setPosition(feature.getGeometry().flatCoordinates);
})


var map = new Map({
  target: 'map',
  layers: [
    new Tile({
      source: new OSM({
        url: `${host}/tile/{z}/{x}/{y}`
      })
    }),
    stationLayer
  ],
  view: new View({
    center: fromLonLat([0.0, 0.0]),
    zoom: 2,
    minZoom: 2,
    maxZoom: 5
  }),
    overlays: [popupOverlay],
});

map.addInteraction(stationSelectInteraction)

/* /* var lineStyle = [new Style({stroke: new Stroke({color: '#d12710', width: 2})})];
 *  *
 *  * var rawPointA = fromLonLat([-5.0, 54.0])
 *  * var rawPointB = fromLonLat([-50.0, 30.0])
 *  * var pointA = new Feature({geometry: new Point(rawPointA)})
 *  * var pointB = new Feature({geometry: new Point(rawPointB)})
 *  *
 *  * var line = new Feature({geometry: new LineString([rawPointA, rawPointB]),
 *  *                            style: lineStyle,
 *  *                            name: "Line :)"})
 *  *
 *  * stationSource.addFeature(pointA)
 *  * stationSource.addFeature(pointB)
 *  * stationSource.addFeature(line) */

search = false
searchBand = null
searchCallsign = null
searchDate = null
searchTime = null

function searchStations(event, options){
  event.preventDefault()

  search = true

  _callsign = document.querySelector('#search input[name="callsign"]').value
  searchCallsign = (_callsign != "" ? _callsign : null)

  _band = document.querySelector('#search input[name="band"]').value
  searchBand = (_band != "" ? _band : null)

  stationFeatures = []
  refreshStations()
  return false
}
document.querySelector('#search').addEventListener("submit", searchStations)

function refreshStations(){

  paramList = []
  params = ""
  if (searchCallsign != null){
    paramList.push(`callsign=${searchCallsign}`)
  }

  if(paramList.length > 0){
    params = "?" + paramList.join('&')
  }

  fetch(`${host}/stations${params}`)
    .then(response => {
      return response.json();
    })
    .then(stations => {

      if (search) stationSource.clear()

      for (let station in stations) {
        if (!(station in stationFeatures) && stations[station]['maidenhead'] != ""){

          stationFeatures[station] = new Feature({
            geometry: new Point(fromLonLat([stations[station]['longitude'], stations[station]['latitude']])),
            data: stations[station]
          })

          stationSource.addFeature(stationFeatures[station]) // this might not be the right way to do this

          console.log(stations[station])
        }
      }

      search = false
    })
}

refreshStations()
setInterval(refreshStations, 15000)

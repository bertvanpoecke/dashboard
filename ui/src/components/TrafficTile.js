import React from "react";
import { Tile } from "./Tile";

export function TrafficTile( props ) {
    const incident = props.incident;
    console.log( incident );
    if ( incident ) {
        return (
            <Tile>
                <span>
                    <div className="road">
                        <span>{incident.road}</span>
                    </div>
                </span>
                <div className="direction">
                    {incident.direction.valueOf()}
                </div>
                <div className="location">
                    <i className="fas fa-map-marker-alt"></i>
                    {incident.location}
                </div>
                <div className="description">
                    <p>{incident.description}</p>
                </div>
            </Tile>
        );
    }
    else {
        return (
            <div>
                ERROR: incident not defined.
            </div>
        );
    }
}






// import React from "react";
// import { Tile } from "./Tile";
// import io from "socket.io-client";

// export default class TrafficTile extends React.Component {
//     constructor(props){
//         super(props);
//         this.state = {
//             trafficData: []
//         };

//         this.socket = io('localhost:3000', {autoConnect: false});
//         this.socket.open();
//         this.socket.on('newTrafficData', function(data){
//             updateData(data);
//         });

//         const updateData = data => {
//             // console.log( data );
//             this.setState( { trafficData: data } );
//         };
//     }

//     render() {
//         const incident = (this.state.trafficData && this.state.trafficData.incidents && this.state.trafficData.incidents.length > 2) ? this.state.trafficData.incidents[2] : null;
//         if ( incident ) {
//             return (
//                 <Tile>
//                     <span>
//                         <div className="road">
//                             <span>{incident.road}</span>
//                         </div>
//                     </span>
//                     <div className="direction">
//                         {incident.direction.valueOf()}
//                     </div>
//                     <div className="location">
//                         <i className="fas fa-map-marker-alt"></i>
//                         {incident.location}
//                     </div>
//                     <div className="description">
//                         <p>{incident.description}</p>
//                     </div>
//                 </Tile>
//             );
//         }
//         else {
//             return (
//                 <div>
//                     niets
//                 </div>
//             );
//         }
//             // <Tile content={this.state.trafficData} />
//     }
// }

import React from "react";
import { TrafficTile } from "./TrafficTile";
// import io from "socket.io-client";
import { CLIENT_RENEG_WINDOW } from "tls";

export default class TrafficSection extends React.Component {
    constructor(props){
        super(props);
        this.state = {
            trafficData: []
        };

        // this.socket = io('localhost:3000', {autoConnect: false});
        // this.socket.open();
        // this.socket.on('newTrafficData', function(data){
        //     updateData(data);
        // });

        const updateData = data => {
            // console.log( data );
            this.setState( { trafficData: data } );
        };
    }

    render() {
        if ( this.state.trafficData && this.state.trafficData.incidents ) {
            const trafficTiles = this.state.trafficData.incidents.map( (incident) =>
                <TrafficTile key={Math.random()} incident={incident}/>
            );
            console.dir( trafficTiles );
            return ( <div>{trafficTiles}</div> );
        }
        // const incident = (this.state.trafficData && this.state.trafficData.incidents && this.state.trafficData.incidents.length > 2) ? this.state.trafficData.incidents[2] : null;
        // if ( incident ) {
        //     return (
        //         <TrafficTile incident={incident}/>
        //     );
        // }
        else {
            return (
                <div>
                    niets
                </div>
            );
        }
    }
}

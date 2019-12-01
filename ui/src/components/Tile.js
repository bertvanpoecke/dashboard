import React from "react";

export function Tile(props) {
    return (
        <div className="tile">
            {props.children}
        </div>
    );
}
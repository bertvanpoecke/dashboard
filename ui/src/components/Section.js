import React from "react";

import Box from '@material-ui/core/Box';
import Container from '@material-ui/core/Container';

export default function Section(props) {
	return (
		<Box m={6}>
			<Container maxWidth="md">
				{props.children}
			</Container>
		</Box>
	)
}
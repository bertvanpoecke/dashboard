import React from "react";

import Section from './Section.js';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import CardMedia from '@material-ui/core/CardMedia';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';

import { withStyles } from '@material-ui/core/styles';

const wikipediaNLURL = "https://nl.wikipedia.org/api/rest_v1/page/random/summary"

const styles = theme => ({
	media: {
		height: 600,
	}
});

class WikiSection extends React.Component {
	constructor(props) {
		super(props);
		this.state = {
			data: null
		};
	}

	componentDidMount() {
		fetch(wikipediaNLURL)
			.then(response => response.json())
			.then(data => this.setState({ data }));
	}

	render() {
		const { classes } = this.props;

		var title = ""
		var summary = ""
		var imagelink = ""
		if (this.state.data != null) {
			title = this.state.data.displaytitle != null ? this.state.data.displaytitle : this.state.data.title
			// summary = this.state.data.extract_html != null ? this.state.data.extract_html : this.state.data.extract
			summary = this.state.data.extract
			imagelink = this.state.data.originalimage != null ? this.state.data.originalimage.source : ""
		}

		return (
			<Section>
				<Card>
					<CardActionArea >
						<CardMedia className={classes.media}
							component="img"
							alt={"Thumbnail " + { title }}
							height="140"
							image={imagelink}
							title={"Thumbnail " + { title }}
						/>
						<CardContent >
							<Typography gutterBottom variant="h5" component="h2">
								{title}
							</Typography>
							<Typography variant="body2" color="textSecondary" component="p">
								{/* <div dangerouslySetInnerHTML={{ __html: summary }} /> */}
								{summary}
							</Typography>
						</CardContent>
					</CardActionArea>
					<CardActions>
						<Button size="small" color="primary">
							Refresh
					</Button>
					</CardActions>
				</Card>
			</Section >
		)
	}
}

export default withStyles(styles)(WikiSection);

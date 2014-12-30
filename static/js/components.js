var Game = React.createClass({
    loadGameFromServer: function() {
        $.ajax({
          url: this.props.url,
          dataType: 'json',
          success: function(data) {
            this.setState({data: data});
          }.bind(this),
          error: function(xhr, status, err) {
            console.error(this.props.url, status, err.toString());
          }.bind(this)
        });
    },
    getInitialState: function() {
        return {
            data: null
        };
    },
    componentDidMount: function() {
        this.loadGameFromServer();
    },
    render: function() {
        if (!this.state.data) {
            return <div className="game"><Loading /></div>;
        }

        var players = this.state.data.result.players;

        console.log(this.state);
        return (
            <div className="game">
                <Player username={players[0].username} /> vs <Player username={players[1].username} />
                <Field />
            </div>
        );
    }
});

var Loading = React.createClass({
    render: function() {
        return (
            <div className="loading">Loading...</div>
        );
    }
});

var Player = React.createClass({
    render: function() {
        return (
            <span>{this.props.username}</span>
        );
    }
})

var Field = React.createClass({
    render: function() {
        var classes = React.addons.classSet({
            "field": true,
        });
        return (
            <div className={classes}>
                <table></table>
            </div>
        );
    }
});

var TurnController = React.createClass({
    render: function() {
        var label = "Waiting";
        if (this.props.turnstate == "pending") { label = "End Turn" };
        return (
            <div className="turn-controller">
                <button className={this.props.turnstate}>{label}</button>
            </div>
        );
    }
});

var Side = React.createClass({
    render: function() {
        var classes = React.addons.classSet({
            "side": true,
            "viewer-side": this.props.owner == "viewer",
            "opponent-side": this.props.owner == "opponent"
        });
        var turncontroller;
        var hand;
        if (this.props.owner == "viewer") {
            turncontroller = <TurnController turnstate={this.props.player.moves ? "pending" : "waiting"} />
            if (this.props.player.hand) {
                hand = <Hand data={this.props.player.hand} />
            };
        };
        return (
            <div className={classes}>
                <Badge name={this.props.player.name} score={this.props.player.score} />
                <Board data={this.props.player.board} />
                {hand}
                {turncontroller}
            </div>
        );
    }
});

var Board = React.createClass({
    render: function() {
        var spaces = this.props.data.map(function (space) {
            return (
                <Space data={space.stack} />
            );
        });
        return (
            <div className="board">
                {spaces}
            </div>
        );
    }
});

var Space = React.createClass({
    render: function() {
        var statbubble;
        var attack = 0;
        var health = 0;
        var tiles = this.props.data.tiles;
        if (tiles) {
            tiles = tiles.map(function (tile) {
                attack += tile.card.attack;
                health += tile.card.health;
                return (
                    <Tile card={tile.card} />
                );
            });
            statbubble = <StatBubble attack={attack} health={health} />;
        };
        return (
            <div className="space">
                {tiles}
                {statbubble}
            </div>
        );
    }
});

var Tile = React.createClass({
    render: function() {
        var imageurl = "/img/tiles/" + this.props.card.name + 'Tile' + imageDensitySuffix + '.png';
        return (
            <div className="tile">
                <img src={imageurl} width="96" height="96" />
                <Card name={this.props.card.name} />
            </div>
        );
    }
});

var StatBubble = React.createClass({
    render: function() {
        return (
            <div className="stat-bubble">
                <span class="attack">{this.props.attack}</span>
                <span class="health">{this.props.health}</span>
            </div>
        );
    }
});

var Hand = React.createClass({
    render: function() {
        var cardCount = this.props.data.length;
        var classes = React.addons.classSet({
            "hand": true,
            "hand-med": cardCount > 3 && cardCount <= 5,
            "hand-large": cardCount > 5 && cardCount <= 7,
            "hand-xlarge": cardCount > 7
        });
        var cards = this.props.data.map(function (card) {
            return (
                <Card name={card.name} />
            );
        });
        return (
            <div className={classes}>
                {cards}
            </div>
        );
    }
});

var Card = React.createClass({
    render: function() {
        var imageurl = "/img/cards/" + this.props.name + 'Card' + imageDensitySuffix + '.png';
        return (
            <div className="card">
                <img src={imageurl} width="160" height="240" />
            </div>
        );
    }
});

var Badge = React.createClass({
    render: function() {
        return (
            <div className="badge">
                <span className="marker">{this.props.score}</span>
                <span className="label">{this.props.name}</span>
            </div>
        );
    }
});

React.render(
    <Game url="/api/game/new?user=Alice&user=Bob" />,
    document.getElementById("main")
);

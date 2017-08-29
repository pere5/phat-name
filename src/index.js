import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';

function Square(props) {
  return (
    <button className="square" onClick={props.onClick}>
      {props.value}
    </button>
  );
}

class Board extends React.Component {

  renderSquare(i) {
    return (
      <Square
        key={'square ' + i}
        value={this.props.squares[i]}
        onClick={() => this.props.onClick(i)}
      />
    );
  }

  render() {
    let index = -1;
    return (
      <div>
        {[...new Array(3)].map((x, i) =>
          <div key={'row ' + i} className="board-row">
            {[...new Array(3)].map((y, j) => {
              index++;
              return this.renderSquare(index);
            })}
          </div>
        )}
      </div>
    );
  }
}

class Game extends React.Component {

  constructor() {
    super();
    this.state = {
      history: [{
        squares: new Array(9).fill(null),
        move: null,
      }],
      xIsNext: true,
      stepNumber: 0,
    };
  }

  handleClick(i) {
    const history = this.state.history.slice(0, this.state.stepNumber + 1);
    const current = history[this.state.stepNumber];
    const squares = current.squares.slice();
    if (calculateWinner(squares) || squares[i]) {
      return;
    }
    squares[i] = this.state.xIsNext ? 'X' : 'O';
    this.setState({
      history: history.concat([{
        squares: squares,
        move: i,
      }]),
      xIsNext: !this.state.xIsNext,
      stepNumber: history.length,
    });
  }

  moveToReadable(move) {
    return move === null ? 'Game start'
        : move === 0 ? 'Move (1, 1)'
        : move === 1 ? 'Move (1, 2)'
        : move === 2 ? 'Move (1, 3)'
        : move === 3 ? 'Move (2, 1)'
        : move === 4 ? 'Move (2, 2)'
        : move === 5 ? 'Move (2, 3)'
        : move === 6 ? 'Move (3, 1)'
        : move === 7 ? 'Move (3, 2)'
        : move === 8 ? 'Move (3, 3)'
        : "";
  }

  jumpTo(step) {
    this.setState({
      stepNumber: step,
      xIsNext: (step % 2) === 0,
    });
  }

  render() {
    const history = this.state.history;
    const current = history[this.state.stepNumber];
    const winner = calculateWinner(current.squares);

    const moves = history.map((board, index) => {
      const desc = this.moveToReadable(history[index].move);
      return (
        <li key={index}>
          <a href="#" onClick={() => this.jumpTo(index)}>{index === this.state.stepNumber ? <b>{desc}</b> : desc}</a>
        </li>
      );
    });

    let status;
    if (winner) {
      status = 'Winner: ' + winner;
    } else {
      status = 'Next player: ' + (this.state.xIsNext ? 'X' : 'O');
    }

    return (
      <div className="game">
        <div className="game-board">
          <Board
            squares={current.squares}
            onClick={(i) => this.handleClick(i)}
          />
        </div>
        <div className="game-info">
          <div>{status}</div>
          <ol>{moves}</ol>
        </div>
      </div>
    );
  }
}

// ========================================

ReactDOM.render(<Game />, document.getElementById('root'));

function calculateWinner(squares) {
  const lines = [
    [0, 1, 2],
    [3, 4, 5],
    [6, 7, 8],
    [0, 3, 6],
    [1, 4, 7],
    [2, 5, 8],
    [0, 4, 8],
    [2, 4, 6],
  ];
  for (let i = 0; i < lines.length; i++) {
    const [a, b, c] = lines[i];
    if (squares[a] && squares[a] === squares[b] && squares[a] === squares[c]) {
      return squares[a];
    }
  }
  return null;
}
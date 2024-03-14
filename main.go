package main

import "time"

type Entry struct {
	// Considered commited if it is safe for that entry to be applied to state machines
	term    int   // Term it was created in
	command State // for state machine?

}
type Cluster struct {
	electionTimeout time.Time // TODO, how long a server waits to do election, 150 ms-300ms

	//Persistent state
	currentTerm int     //latest term server has seen (initialized to 0 on first boot, increases monotonically)
	votedFor    int     //candidateId that received vote in current term (or null if none)
	log         []Entry //log entries; each entry contains command for state machine, and term when entry was received by leader (first index is 1)

	//Volatile state
	commitIndex int //index of highest log entry known to be committed (initialized to 0, increases monotonically)
	lastApplied int //index of highest log entry applied to state machine (initialized to 0, increases monotonically)
}
type State string

/*
- Followers only respond to requests from leaders
- Leaders handles client requests
- Candidate is used to elect a new leader
*/
const (
	Leader    State = "leader"
	Follower        = "follower"
	Candidate       = "candidate"
)

type Node struct {
	state State
	// Volatile state on Leaders
	nextIndex  []int //for each server, index of the next log entry to send to that server (initialized to leader last log index + 1)
	matchIndex []int //for each server, index of highest log entry known to be replicated on server (initialized to 0, increases monotonically)
}

// AppendEntriesRequest Initiated by the leader to replicate log entries
// AppendEntriesRequest  Heartbeats are appendentries without entries
type AppendEntriesRequest struct {
	term         int     //leader’s term
	leaderId     int     //so follower can redirect clients
	prevLogIndex int     //index of log entry immediately preceding new ones
	prevLogTerm  Entry   //term of prevLogIndex entry
	entries      []Entry //log entries to store (empty for heartbeat; may send more than one for efficiency)

	leaderCommit int //leader’s commitIndex
}

type AppendEntriesResult struct {
	term    int  //currentTerm, for leader to update itself
	success bool //true if follower contained entry matching prevLogIndex and prevLogTerm
}
type RequestVoteRequest struct {
	term         int //candidate’s term
	candidateId  int //candidate requesting vote
	lastLogIndex int //index of candidate’s last log entry
	lastLogTerm  int //term of candidate’s last log entry
}

type RequestVoteResult struct {
	term        int  //currentTerm, for candidate to update itself
	voteGranted bool //true means candidate received vote
}

func LeaderElection() {

	/*
		- Heartbeat is used to trigger an election
		- Servers startup as followers, until it receives a valid RPC from leader or candidate
		- Leaders send periodic heartbeats, append entries rpc without entries to all followers to maintain authority
		- If follower receives none for a certain time, called election timeout, goes through election

		Begin Election
			- follower increments current term, transitions to candidate state
			- Vote for itself, issues a RequestVote RPC in parallel to the other servers in the cluster
			- Remains in this state until
				- a) wins the election
				- b) another server estabilishes itself as leader
				- c) period of time goes by without a leader
			- Voted on first come first serve
			- When one wins, a heartbeat is sent out to all the servers, to indicate elections r done


			-When waiting for votes, candidate may get AppendEntries RPC. if the leaders term, in the rpc,
			is atleast as large as the candidates current term, then it sees it as legitamate, and returns to follower state
			Otherwise, reject and continue as candidate state

			- Neither win/lose, each candidate will timeout, start new election, by incrementing its term and another round of Request vote rpcs


	*/
}

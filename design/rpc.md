TCPを使ってノード間のRequest/Responseを送受信する。

#### Codec

- `encode(message: Message) -> []Bytes`
- `decode(payload: Bytes) -> Message`

#### Framing

- `writeFrame(stream, payload)`
- `readFrame(stream) -> []Bytes`


Format:
```
[4-byte payload length][1-byte message type][N-byte payload]
```


#### Client/Server

- `call(peer, request) -> Response`
- `serve(listener, handler)`

### Message

#### Propose

ProposeRequest/ProposeResponse: クライアントから送信されるメッセージとそれに対する返信

ProposeRequest { data: bytes }
ProposeResponse { status, term, logindex, leaderId } # リーダーだけに送られるわけではないので、それをハンドリングしないといけない。

Status:

indexやnodeIdのzero valueはunknownを示すようにする。

- Commited
- Failed ( term: 応答ノードのcurrent Term, logIndex: 0, leaderId: ノードID )
- NotLeader ( term: 応答ノードのcurrent term, logIndex: 0, leaderId: 0 )

#### RequestVote

RequestVoteRequest { Term, CandidateID, LatestLogIndex, LatestLogTerm }
RequestVoteResponse { VoterTerm, VoterID, VoteGranted } 

VoteをRequestしたノードは、Response.VoterTermをみて、もし自分の方がtermが小さい場合、followerに降格しないといけない。

またネットワークリトライなどによってVoteリクエストが再度届いた際には、ファイルに永続化されたvotedForを使って冪等性を担保
したまま処理ができる。

#### AppendEntries

LogEntry {Index, Term, Payload}
AppendEntriesRequest { Term, LeaderID, PrevLogIndex, PrevLogTerm, Entries, LeaderCommit }

PrevLogIndex, PrevLogTerm: followerとのログ整合性チェックで使う。
LeaderCommit: Leaderがコミット済みと判断しているログのindex。ここまでシステム全体でcommit済みであることを
示すwatermark. Followerからみると、エントリを受け取ったがcommitされているかはわからない。
HeartbeatでもLeaderCommitを送るようにして、ログの保存とcommitの認知を後から知るようになる。そうしないと、
読み取り時に、どのindexまでが有効かを判断できない。

AppendEntriesResponse { Term, FollowerID, Success, MatchIndex }

Status:

- Success
- Failed


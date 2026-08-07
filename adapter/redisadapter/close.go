package redisadapter

func (a Adapter) Close() error {

	return a.client.Close()

}

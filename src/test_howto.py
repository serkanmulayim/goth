from balta import *


api_key="Oos1jUZpXSEmIxxTF5V0rFMmMRpSsopzWhCjobiLPkyzpJ36WDHYwlwzGPWfEskJ";
api_secret="kE5ORsdAyyA6wNdm85fiHPwhCt51iLsQWXkxvqR71lI6maf9AFEB4H3C7rJZ8uuG";
DATA_FOLDER = "../data/";

symbol = "BTCUSDT"
interval = Client.KLINE_INTERVAL_1DAY

start =  "August 01, 2017"
end = "February 01, 2018"

#get data from binance
balta = Balta() #no parameter for reading from binance
balta.get_data_from_exchange('binance',api_key, api_secret, symbol, interval, start, end)
balta.write_data_to_csv('../data/binance.csv')

#read data from file
df = pd.read_csv('../data/binance.csv', index_col=0,
	converters={'open':Decimal, 'high':Decimal, 'low':Decimal, 'close':Decimal, 'volume':Decimal, 'qav':Decimal, 'numtrades':Decimal,'tbbav':Decimal, 'tbqav':Decimal}
	)

balta2 = Balta(df);


sma10 = EMA(20);
balta2.add_indicator_data(sma10);
df.to_csv("../data/data2.csv")
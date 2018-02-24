import numpy as np
import pandas as pd
from indicators import *
from binance.client import Client
import os

COL_NAME_OPEN = 'open'
COL_NAME_HIGH = 'high'
COL_NAME_LOW = 'low'
COL_NAME_CLOSE = 'close'
COL_NAME_VOLUME = 'volume'
COL_NAME_QUOTE_ASSET_VOLUME = 'qav'
COL_NAME_NUM_TRADES = 'numtrades'
COL_NAME_TAKER_BUY_BASE_ASSET_VOLUME = 'tbbav'
COL_NAME_TAKER_BUY_QUOTE_ASSET_VOLUME = 'tbqav'

BALTA_INITIAL_COLS = cols = [COL_NAME_OPEN, COL_NAME_HIGH, COL_NAME_LOW, COL_NAME_CLOSE, COL_NAME_VOLUME, COL_NAME_QUOTE_ASSET_VOLUME, COL_NAME_NUM_TRADES, COL_NAME_TAKER_BUY_BASE_ASSET_VOLUME, COL_NAME_TAKER_BUY_BASE_ASSET_VOLUME]

class Balta(object):



    def __init__(self, dataframe=None):
        self.dataframe = dataframe        

    def add_indicator_data(self, indicator):
    	cols = self.get_indicator_cols(indicator)
    	if len(cols) == 0:
    		print ("Unknown indicator " + indicator.label +". Nothing is added to the data")
    		return

    	data = np.array(self.dataframe[cols])    	
    	data_to_add = np.squeeze(np.array(indicator.calculate(data)))
    	#add new col to the dataframe
    	self.dataframe[indicator.label] = data_to_add


    def get_indicator_cols(self, indicator):
    	if indicator.label.startswith('SMA') or indicator.label.startswith('EMA') or indicator.label.startswith('MACD'):
    		return ['close']
    	return []

    def get_data_from_exchange(self, exchange, api_key, api_secret, symbol, interval, start_date, end_date):
    	#from any exchange to data
    	if exchange == 'binance':
    	    client = Client(api_key, api_secret)
    	    res = client.get_historical_klines(symbol, interval, start_date, end_date)
    	    np_res = np.array(res)
    	    floats = np_res.astype(np.float)
    	    index = np.transpose(floats)[0].astype(np.int)
            data = np.delete(floats,0, axis=1) #delete opentime col
            data = np.delete(data,5, axis=1) #delete closetime col
            data = np.delete(data,9, axis=1) #can be ingnored
            self.dataframe = pd.DataFrame(data, index= index, columns = BALTA_INITIAL_COLS)
    	else:
    		raise Exception('unknown exchange:' + exchange)

    def write_data_to_csv(self, file):
        if self.dataframe is None:
        	raise Exception('There is no dataframe initialized. No file is created')
        dirr = os.path.dirname(file)
        if not os.path.exists(dirr):
            try:
                os.makedirs(dirr)
            except OSError as exc:
                raise Exception('error occured while creating the path:' + dirr)
        self.dataframe.to_csv(file)
        







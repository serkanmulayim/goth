import numpy as np
import pandas as pd
from indicators import *

class Balta(object):
    def __init__(self, dataframe):
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




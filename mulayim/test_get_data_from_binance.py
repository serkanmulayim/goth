#"KTIDRUSEAP4BYFYL" // binance Google Auth recovery key

from binance.client import Client

import pandas as pd
import numpy as np
from balta import *

api_key="Oos1jUZpXSEmIxxTF5V0rFMmMRpSsopzWhCjobiLPkyzpJ36WDHYwlwzGPWfEskJ";
api_secret="kE5ORsdAyyA6wNdm85fiHPwhCt51iLsQWXkxvqR71lI6maf9AFEB4H3C7rJZ8uuG";
DATA_FOLDER = "../data/";

symbol = "BTCUSDT"
interval = Client.KLINE_INTERVAL_1DAY

start =  "August 01, 2017"
end = "February 01, 2018"

balta = Balta()

balta.get_data_from_exchange('binance',api_key, api_secret, symbol, interval, start, end)
balta.write_data_to_csv('../data/binance.csv')


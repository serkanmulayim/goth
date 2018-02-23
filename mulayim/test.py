import numpy as np
import pandas as pd
from indicators import *
from balta import *
from decimal import Decimal

df = pd.read_csv('../data/data.csv', index_col=0,
	converters={'open':Decimal, 'high':Decimal, 'low':Decimal, 'close':Decimal, 'volume':Decimal, 'qav':Decimal, 'numtrades':Decimal,'tbbav':Decimal, 'tbqav':Decimal}
	)

balta = Balta(df);
sma10 = EMA(20);
balta.add_indicator_data(sma10);
df.to_csv("../data/data2.csv")
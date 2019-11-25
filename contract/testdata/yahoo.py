
# In[69]:


import pandas as pd
import yfinance as yf
import datetime

aapl = yf.Ticker("AAPL")
flat_df = pd.DataFrame()
for expi in aapl.options:
    chain = aapl.option_chain(expi)
    chain.calls["putcall"] = "call"
    chain.puts["putcall"] = "put"

    flat_df = flat_df.append(chain.calls)
    flat_df = flat_df.append(chain.puts)
    flat_df["expiry"] = datetime.datetime.strptime(expi, '%Y-%m-%d')

flat_df["undticker"] = "AAPL"
flat_df = flat_df[["contractSymbol", "strike", "undticker", "expiry", "putcall"]]
flat_df = flat_df.rename(columns={'contractSymbol': 'ticker'})

print flat_df.head(10)


# In[70]:


flat_df.to_json(os.path.join(os.getcwd(),"yahoo.json"), orient='records', date_format='epoch')


# In[ ]:


import simplejson as json
import requests

contracts = json.loads(flat_df.to_json(orient='records', date_format='epoch'))
res = [requests.post('http://localhost:3000/europeans', json=contract) for contract in contracts]

print res.head(10)


# In[ ]:


res = requests.get('http://localhost:3000/europeans')
europeans = res.json()['europeans']

print json.dumps(europeans[:20], indent=4)



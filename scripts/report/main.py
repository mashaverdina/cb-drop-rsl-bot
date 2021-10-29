import pandas as pd
import numpy as np
import plotly.express as px
import plotly.graph_objs as go
import scipy

DROP_ITEMS = ['ancient_shard', 'void_shard', 'sacred_shard', 'epic_tome', 'leg_tome']
MAPPING = {
    'ancient_shard': 'Древний',
    'void_shard': 'Войд',
    'sacred_shard': 'Сакрал',
    'epic_tome': 'Эпик том',
    'leg_tome': 'Лег том',
}
USER = 'user_id'
BOSS = 'level'
RELATED_TO = 'related_to'

CONFIG = {
    'days_low_bound': 25,
    'bootstrap_count': 100,
}

class DayStat:
    def __init__(self, label, arr, trust_intervals, index):
        self.arr = arr
        self.size = len(arr)
        self.trust_intervals = trust_intervals
        self.label = label
        self.index = index

    def plot_data(self):

        trace = go.Scatter(
            y=self.arr,
            x=self.index,
            name=self.label,
            error_y=dict(
                type='data', # value of error bar given in data coordinates
                # symmetric=False,
                array=[x[1] - x[0] for x in self.trust_intervals],
                # arrayminus=[x[0] for x in self.trust_intervals],
                visible=True,
            )
        )
        return [trace], None, None, None

class GeneralStat:
    def __init__(self, label, arr):
        arr = 30 * arr
        self.mean = np.mean(arr)
        self.p25 = np.quantile(arr, 0.25)
        self.p50 = np.quantile(arr, 0.5)
        self.p75 = np.quantile(arr, 0.75)
        self.p98 = np.quantile(arr, 0.98)
        self.p99 = np.quantile(arr, 0.99)
        self.arr = arr
        self.size = len(arr)
        self.label = label

    def plot_data(self, nbinsx=15):
        nbinsx = int(max(self.arr) - min(self.arr))
        hist = go.Histogram(
            x=self.arr,
            opacity=0.75,
            histnorm='probability',
            name=self.label,
            nbinsx=nbinsx,
        )

        shape = {'line':
                     {'color': '#0099FF', 'dash': 'solid', 'width': 1},
                     'type': 'line',
                     'x0': self.p50,
                     'x1': self.p50,
                     'xref': 'x',
                     'y0': -0.1,
                     'y1': 1,
                     'yref': 'paper'
                 }

        annotation = dict(
                x=self.p50,
                y=1,
                xref='x',
                yref='paper',
                text=f'Med {self.p50:2.2f}',
                showarrow=True,
                arrowhead=7,
                ax=1,
                ay=1,
                # axref='paper',
                # ayref='paper'
            )

        return None, hist, shape, annotation


# Index(['user_id', 'related_to', 'level', 'last_update', ],
def load_data(path='data/dump29.csv'):
    df = pd.read_csv(path, sep=',')
    df = df[df[RELATED_TO] >= '2021-10-00 00:00:00']
    # df['related_to'] = pd.to_datetime(df['related_to'], format="%Y-%m-%d")
    return df

def clean(df):
    for item in DROP_ITEMS:
        df = df[df[item] <= 2]
    return df


def by_user_stat(df):
    aggregations = {item: ['sum', 'mean'] for item in DROP_ITEMS}
    aggregations[USER] = 'count'
    aggregated = df.groupby([BOSS, USER]) \
        .agg(aggregations)
    aggregated['days'] = aggregated[USER, 'count']
    # aggregated.d
    return aggregated

def trust_interval(series):
    repl = np.array([np.mean(x) for x in (np.random.choice(series.values, len(series)) for _ in range(CONFIG['bootstrap_count']))])
    return (np.quantile(repl, 0.05), np.quantile(repl, 0.95))

def by_day_stat(df):
    aggregations = {item: ['sum', 'mean', trust_interval] for item in DROP_ITEMS}
    # aggregations[USER] = 'count'
    aggregated = df.groupby([BOSS, RELATED_TO]) \
        .agg(aggregations)
    # aggregated['days'] = aggregated[USER, 'count']
    # aggregated.d
    return aggregated


def filter_by_days(df, days=None):
    if not days:
        days = CONFIG['days_low_bound']
    return df[df['days'] > days]
    
def create_fig(title, stats):
    hists = []
    shapes = []
    annotations = []
    all_traces = []
    for stat in stats:
        traces, hist, shape, annotation = stat.plot_data()
        if hist:
            hists.append(hist)
        if traces:
            all_traces.extend(traces)
        if shape:
            shapes.append(shape)
        if annotation:
            annotations.append(annotation)

    layout = go.Layout(
        title=title,
        barmode='overlay',
        xaxis=dict(
            title='Айтемов в день'
        ),
        yaxis=dict(
            title='Доля пользователей'
        ),
        shapes = shapes,
        annotations=annotations,
    )
    if len(hists):
        fig = go.Figure(data=hists, layout=layout)
    else:
        fig = go.Figure(layout=layout)
    for trace in all_traces:
        fig.add_trace(trace)
    return fig

if __name__ == '__main__':
    df = load_data()
    print(f'loaded {len(df)} items')
    df = clean(df)
    print(f'after clean: {len(df)}')
    # print(df.head())
    # print(df.columns)
    print('--------------')
    udf = by_user_stat(df)
    print(f'got {len(udf)} per user records')
    udf = filter_by_days(udf)
    print(f'got {len(udf)} per user records after days filter')
    ddf = by_day_stat(df)
    print(f'got {len(ddf)} per days records')

    print(ddf.head())
    # print(udf.head())
    # print(udf.columns)
    # # print(udf[157683409, 6])
    # print('------------------------------------------')
    # # print(udf.loc[[157683409],:])
    # print(udf.loc[(5, 157683409),:])
    # print('------------------------------------------')
    # print(udf.loc[(6, 157683409),:])
    # print(type(udf['ancient_shard']['mean'].values))

    for boss in [5, 6]:
        generalStats = []
        for item in ['ancient_shard', 'void_shard', 'sacred_shard']:
            x = GeneralStat(f'{boss}КБ {MAPPING[item]}', udf.loc[boss][item]['mean'].values)
            generalStats.append(x)
            fig = create_fig(f'Распределение дропа осколков {boss}КБ {item}', [x])
            fig.show()
        fig = create_fig(f'Распределение дропа осколков {boss}КБ', generalStats)
        fig.show()

        generalStats = []
        for item in ['epic_tome', 'leg_tome']:
            generalStats.append(GeneralStat(f'{boss}КБ {MAPPING[item]}', udf.loc[boss][item]['mean'].values))
        fig = create_fig(f'Распределение дропа книжек {boss}КБ', generalStats)
        fig.show()

    for boss in [5, 6]:
    # for boss in [6]:
        dayStats = []
        # for item in (x for x in DROP_ITEMS if not (x == 'ancient_shard' and boss == 6)):
        for item in DROP_ITEMS:
        # for item in ('ancient_shard'):
            data = ddf.loc[boss][item]

            dayStats.append(DayStat(f'{boss}КБ {MAPPING[item]}', data['mean'].values, data['trust_interval'].values, data['mean'].index.values))
        fig = create_fig(f'Средний дроп по дням {boss}КБ', dayStats)
        fig.show()
    # print(boss5.loc[157683409])
    # print(udf.loc[[157683409],:].hgrad())
# See PyCharm help at https://www.jetbrains.com/help/pycharm/

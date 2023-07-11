const getTimestampCount = (count: string) =>
  parseInt((Date.now() / 1000).toString()) + "-" + count;

export { getTimestampCount };

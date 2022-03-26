import { isTrue, normalizeStringValue, trim } from './tools';

describe('Tools test', () => {
  it('Should collect true value from org string', () => {
    expect(isTrue('yes')).toEqual(true);
    expect(isTrue('true')).toEqual(true);
    expect(isTrue('    true   ')).toEqual(true);
    expect(isTrue('    yes  ')).toEqual(true);
  });

  it('Should collect false value from org string', () => {
    expect(isTrue('yes some text')).toEqual(false);
    expect(isTrue('tetrue')).toEqual(false);
    expect(isTrue('')).toEqual(false);
    expect(isTrue('1')).toEqual(false);
  });

  it('Should normalize org text', () => {
    expect(normalizeStringValue('   some text with BIG WORDs   ')).toEqual('some text with big words');
  });

  it('Should not change normal string after normalization', () => {
    expect(normalizeStringValue('text should not be normalized')).toEqual('text should not be normalized');
  });

  it('Should normalize upper case text', () => {
    expect(normalizeStringValue('TEXT SHOULD NOT BE NORMALIZED')).toEqual('text should not be normalized');
  });

  it('Should preserve empty string after normalization', () => {
    expect(normalizeStringValue('')).toEqual('');
  });
});

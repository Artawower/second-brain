import { OrgData, OrgNode } from 'uniorg';
import { Note } from './models';
export declare type NodeMiddleware = (orgData: OrgNode) => OrgNode;
export declare const collectNote: (content: OrgData, middlewareChains?: NodeMiddleware[]) => Note;
